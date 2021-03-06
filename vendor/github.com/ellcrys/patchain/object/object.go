package object

import (
	"fmt"
	"strings"
	"time"

	"github.com/ellcrys/cocoon/core/common"
	"github.com/ellcrys/patchain"
	"github.com/ellcrys/patchain/cockroach/tables"
	"github.com/ellcrys/util"
	"github.com/jinzhu/copier"
	"github.com/ncodes/redo"
	"github.com/pkg/errors"
)

// Object defines a structure for handling objects
type Object struct {
	db patchain.DB
}

// NewObject creates a new object handler
func NewObject(db patchain.DB) *Object {
	return &Object{db: db}
}

// Create creates an object to represent anything or resource
func (o *Object) Create(obj *tables.Object) error {
	return o.db.Create(obj.Init().ComputeHash())
}

// CreateOnce creates the object only if no other object shares the same key
func (o *Object) CreateOnce(obj *tables.Object) error {
	var existing tables.Object
	err := o.db.GetLast(&tables.Object{Key: obj.Key}, &existing)
	if err != nil {
		if common.CompareErr(err, patchain.ErrNotFound) == 0 {
			err = o.Create(obj)
		}
		return err
	}
	copier.Copy(obj, existing)
	return nil
}

// CreatePartitions creates partitions. Every partition is chained to the
// one before it by sharing the hash of the previous partition as the new
// partition's prev hash value.
func (o *Object) CreatePartitions(n int64, ownerID, creatorID string, options ...patchain.Option) ([]*tables.Object, error) {

	// process options
	dbTx := o.db.Begin()
	dbOptions := []patchain.Option{&patchain.UseDBOption{DB: dbTx}}
	finish := true
	if len(options) > 0 {
		dbOptions = options
		for _, ops := range options {
			if ops.GetName() == patchain.UseDBOptionName {
				dbOpt := ops.(*patchain.UseDBOption)
				dbTx = dbOpt.GetValue().(patchain.DB)
				finish = dbOpt.Finish
			}
		}
	}

	// Process starts a transaction to create the partition(s). It can be called multiple times when
	// a retry error is returned or a unique constraint error occurs on the prev_hash field
	var process = func() ([]*tables.Object, error) {

		partitions := []*tables.Object{}

		// create partitions
		for i := int64(0); i < n; i++ {
			newPartition := MakePartitionObject(util.UUID4(), ownerID, creatorID)
			partitions = append(partitions, newPartition)
		}

		// chain partitions with prefixed prev hash
		MakeChain(partitions...)

		return partitions, o.db.TransactWithDB(dbTx, finish, func(dbTx patchain.DB, commit patchain.CommitFunc, rollback patchain.RollbackFunc) error {

			// get the last partition
			lastPartition, err := o.GetLast(&tables.Object{QueryParams: patchain.KeyStartsWith(PartitionPrefix)}, dbOptions...)
			if err != nil {
				if err != patchain.ErrNotFound {
					return errors.Wrap(err, "failed to get last partition")
				}
			}

			// if no last partition, add the new partitions
			partitionsI, _ := util.ToSliceInterface(partitions)
			if lastPartition == nil {

				if err := o.db.CreateBulk(partitionsI, dbOptions...); err != nil {
					return errors.Wrap(err, "failed to create partition")
				}

				// create genesis pair for each partition
				for _, partition := range partitions {
					gPair := MakeGenesisPair(ownerID, creatorID, partition.ID, partition.Hash)
					gPairI, _ := util.ToSliceInterface(gPair)
					if err := o.db.CreateBulk(gPairI, dbOptions...); err != nil {
						return errors.Wrap(err, "failed to create genesis object")
					}
				}

				return nil
			}

			// update the previous hash of the first of the new partitions
			// to the hash of the last partition
			partitionsI[0].(*tables.Object).PrevHash = lastPartition.Hash
			MakeChain(partitions...)
			if err := o.db.CreateBulk(partitionsI, dbOptions...); err != nil {
				return errors.Wrap(err, "failed to create partition")
			}

			// create genesis pair for each partition
			for _, partition := range partitions {
				gPair := MakeGenesisPair(ownerID, creatorID, partition.ID, partition.Hash)
				gPairI, _ := util.ToSliceInterface(gPair)
				if err := o.db.CreateBulk(gPairI, dbOptions...); err != nil {
					return errors.Wrap(err, "failed to create genesis object")
				}
			}

			return nil
		})
	}

	partitions, err := process()

	return partitions, errors.Wrap(err, "failed to create partition(s)")
}

// Retry runs an operation if it fails dues to a retry or prev_hash contention error
func (o *Object) Retry(cb func(stop func()) error) error {
	var err error
	c := redo.NewDefaultBackoffConfig()
	c.MaxElapsedTime = 10 * time.Minute
	err = redo.NewRedo().BackOff(c, func(stop func()) error {
		err = cb(stop)
		if err != nil && o.RequiresRetry(err) {
			return err
		}
		stop()
		return err
	})
	return err
}

// MustCreatePartitions is the same as CreatePartitions but it will retry the operation
// if it fails because of a transaction retry or contention error. The retry will continue
// to happen for max of 10 minutes using an exponential backoff algorithm.
// Note: if an external database connection is passed using the db option and a transaction
// fails, it will not be retried
func (o *Object) MustCreatePartitions(n int64, ownerID, creatorID string, options ...patchain.Option) ([]*tables.Object, error) {
	var partitions []*tables.Object
	var err error
	err = o.Retry(func(stop func()) error {
		partitions, err = o.CreatePartitions(n, ownerID, creatorID, options...)
		if err != nil {
			return err
		}
		return nil
	})
	return partitions, err
}

// GetLast gets the latest version of an object.
// It does this by enforcing a descending order of the insert timestamp of the object.
func (o *Object) GetLast(q patchain.Query, options ...patchain.Option) (*tables.Object, error) {
	var obj tables.Object
	err := o.db.GetLast(q, &obj, options...)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

// All fetches all the objects matching a query
func (o *Object) All(q patchain.Query, options ...patchain.Option) ([]*tables.Object, error) {
	var objs []*tables.Object
	err := o.db.GetAll(q, &objs, options...)
	return objs, err
}

// given a slice of partitions, it randomly chooses a partition if
// there are more than one partition, otherwise it returns the only
// partition if only one exists or nil if not exist
func (o *Object) selectPartition(partitions []*tables.Object) *tables.Object {
	if len(partitions) == 0 {
		return nil
	} else if len(partitions) == 1 {
		return partitions[0]
	} else {
		return partitions[util.RandNum(0, len(partitions)-1)]
	}
}

// RequiresRetry checks whether a transaction error
// indicates or requires a retry. This method can detect cockroach db
// restart, retry error and prev hash contention
func (o *Object) RequiresRetry(err error) bool {
	return strings.Contains(err.Error(), "restart transaction") || strings.Contains(err.Error(), "retry transaction") || strings.Contains(err.Error(), `violates unique constraint "idx_prev_hash"`)
}

// MustPut is the same as Put but it will retry the operation if it
// fails because of a transaction retry or contention error.
// The retry will continue to occur for max of 10 minutes using
// an exponential backoff algorithm.
// Note: if an external database connection is passed using the db option and a transaction
// fails, it will not be retried
func (o *Object) MustPut(objs interface{}, options ...patchain.Option) error {
	var err error
	err = o.Retry(func(stop func()) error {
		return o.Put(objs, options...)
	})
	return err
}

// Put adds an object into a randomly selected partition belonging to
// the owner of the object. If object has no owner, error is returned
func (o *Object) Put(objs interface{}, options ...patchain.Option) error {

	var objects []*tables.Object
	switch o := objs.(type) {
	case []*tables.Object:
		objects = o
	case *tables.Object:
		objects = append(objects, o)
	default:
		return fmt.Errorf("unsupported object type")
	}

	if len(objects) == 0 {
		return fmt.Errorf("no object to put")
	}

	// ensure all objects have an owner id and all objects must belong to same owner
	// we need the owner id to determine the partition to add the object to.
	ownerID := objects[0].OwnerID
	for i, o := range objects {
		if o.OwnerID == "" {
			return fmt.Errorf("object %d: object does not have an owner", i)
		} else if o.OwnerID != ownerID {
			return fmt.Errorf("object %d: has a different owner", i)
		}
	}

	// process options
	dbTx := o.db.Begin()
	dbOptions := []patchain.Option{&patchain.UseDBOption{DB: dbTx}}
	finish := true
	if len(options) > 0 {
		dbOptions = options
		for _, ops := range options {
			if ops.GetName() == patchain.UseDBOptionName {
				dbTx = ops.(*patchain.UseDBOption).GetValue().(patchain.DB)
				finish = ops.(*patchain.UseDBOption).Finish
			}
		}
	}

	// define function to perform put operation. May be repeated if the following conditions occur:
	// - Error indicating a restart or retry the transaction
	// - Error indicating a previous hash unique index violation
	var putTxFunc = func() error {

		return o.db.TransactWithDB(dbTx, finish, func(dbTx patchain.DB, commit patchain.CommitFunc, rollback patchain.RollbackFunc) error {

			// get the partitions belonging to the owner of the object
			partitions, err := o.All(&tables.Object{OwnerID: ownerID, QueryParams: patchain.KeyStartsWith(PartitionPrefix)}, dbOptions...)
			if err != nil {
				return errors.Wrap(err, "failed to get owner's partition")
			}

			// select a random partition
			var selectedPartition = o.selectPartition(partitions)
			if selectedPartition == nil {
				return fmt.Errorf("owner has no partition")
			}

			// assign selected partition to the objects
			for _, o := range objects {
				o.PartitionID = selectedPartition.ID
			}

			// get the last object of the selected partition
			lastObj, err := o.GetLast(&tables.Object{PartitionID: selectedPartition.ID}, dbOptions...)
			if err != nil {
				// no object in this partition! This means no genesis pair/object, return error
				if err == patchain.ErrNotFound {
					return fmt.Errorf("no genesis object in the partition")
				}
				return err
			}

			// assign hash of last object as the PrevHash value
			// of the first object, chain the  objects and create them
			objects[0].PrevHash = lastObj.Hash
			MakeChain(objects...)
			for _, o := range objects {
				if err := dbTx.Create(o, options...); err != nil {
					return errors.Wrap(err, "failed to add object to partition")
				}
			}

			// update peer hash of last object
			lastObj.ComputePeerHash(objects[0].Hash)
			if err := dbTx.UpdatePeerHash(lastObj, lastObj.PeerHash, options...); err != nil {
				return errors.Wrap(err, "failed to update last object peer hash")
			}

			return nil
		})
	}

	return errors.Wrap(putTxFunc(), "failed to put object(s)")
}
