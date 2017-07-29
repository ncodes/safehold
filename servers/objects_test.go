package servers

import (
	"context"
	"testing"
	"time"

	"github.com/ellcrys/elldb/servers/db"
	"github.com/ellcrys/elldb/servers/oauth"
	"github.com/ellcrys/elldb/servers/proto_rpc"
	"github.com/ellcrys/util"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/metadata"
)

func TestObject(t *testing.T) {
	setup(t, func(rpc, rpc2 *RPC) {
		Convey("Object", t, func() {
			c1 := &proto_rpc.CreateAccountMsg{FirstName: "john", LastName: "Doe", Email: util.RandString(5) + "@example.com", Password: "something", Developer: true}
			resp, err := rpc.CreateAccount(context.Background(), c1)
			So(err, ShouldBeNil)
			var account map[string]interface{}
			util.FromJSON(resp.Object, &account)

			c2 := &proto_rpc.CreateAccountMsg{FirstName: "john2", LastName: "Doe2", Email: util.RandString(5) + "@example.com", Password: "something", Developer: true}
			resp, err = rpc.CreateAccount(context.Background(), c2)
			So(err, ShouldBeNil)
			var account2 map[string]interface{}
			util.FromJSON(resp.Object, &account2)

			ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
			b := &proto_rpc.CreateBucketMsg{Name: util.RandString(5)}
			bucket, err := rpc.CreateBucket(ctx, b)
			So(err, ShouldBeNil)
			So(bucket.Name, ShouldEqual, b.Name)
			So(bucket.ID, ShouldHaveLength, 36)

			b2 := &proto_rpc.CreateBucketMsg{Name: util.RandString(5), Immutable: true}
			_, err = rpc.CreateBucket(ctx, b2)
			So(err, ShouldBeNil)

			mapping := &proto_rpc.CreateMappingMsg{
				Name:   util.RandString(5),
				Bucket: b.Name,
				Mapping: util.MustStringify(map[string]interface{}{
					"name": "key",
				}),
			}
			_, err = rpc.CreateMapping(ctx, mapping)
			So(err, ShouldBeNil)

			Convey(".CreateObjects", func() {

				Convey("Should return error if bucket is not provided", func() {
					_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: ""})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "400",
						"field":   "bucket",
						"message": "bucket name is required",
					})
				})

				Convey("With bucket", func() {

					Convey("Should return error if bucket does not exists", func() {
						_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: "unknown"})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"field":   "bucket",
							"message": "bucket not found",
							"status":  "404",
						})
					})

					Convey("Should return error if validation fails", func() {
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						objs := util.MustStringify([]map[string]interface{}{
							{"owner_id": 123, "key": "mykey"},
						})
						_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"code":    "invalid_parameter",
							"message": "object 0: owner_id must be a string",
						})
					})

					Convey("Should return error if owner of object does not exist", func() {
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						objs := util.MustStringify([]map[string]interface{}{
							{"owner_id": "unknown", "key": "mykey"},
						})
						_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"code":    "invalid_parameter",
							"message": "owner of object(s) does not exist",
						})
					})

					Convey("Should return permission error if authenticated account does not have permission to create object for the owner", func() {
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						objs := util.MustStringify([]map[string]interface{}{
							{"owner_id": account2["id"].(string), "key": "mykey"},
						})
						_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"status":  "401",
							"message": "permission denied: you are not authorized to create objects for the owner",
						})
					})

					Convey("Should successfully create new objects", func() {
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						objs := util.MustStringify([]map[string]interface{}{
							{"owner_id": account["id"].(string), "key": "mykey", "value": "myval"},
							{"owner_id": account["id"].(string), "key": "mykey2", "value": "myval2"},
						})
						resp, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs})
						So(err, ShouldBeNil)
						var m []map[string]interface{}
						err = util.FromJSON(resp.Objects, &m)
						So(err, ShouldBeNil)
						So(m, ShouldHaveLength, 2)
						So(m[0]["bucket"], ShouldEqual, b.Name)
						So(m[1]["bucket"], ShouldEqual, b.Name)
						So(m[0]["creator_id"], ShouldEqual, account["id"])
						So(m[1]["creator_id"], ShouldEqual, account["id"])
						So(m[0]["owner_id"], ShouldEqual, account["id"])
						So(m[1]["owner_id"], ShouldEqual, account["id"])
					})

					Convey("With mapping", func() {
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						mappingName := util.RandString(5)
						mapping := `{ "user_id": "owner_id", "first_name": "key", "last_name": "value" }`
						_, err := rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{Name: mappingName, Bucket: b.Name, Mapping: []byte(mapping)})
						So(err, ShouldBeNil)

						Convey("Should return error if mapping does not exists", func() {
							ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
							objs := util.MustStringify([]map[string]interface{}{
								{"owner_id": account["id"].(string), "key": "mykey", "value": "myval"},
							})
							_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs, Mapping: "unknown"})
							So(err, ShouldNotBeNil)
							m, err := util.JSONToMap(err.Error())
							So(err, ShouldBeNil)
							errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
							So(errs, ShouldHaveLength, 1)
							So(errs[0], ShouldResemble, map[string]interface{}{
								"message": "mapping not found",
								"status":  "400",
							})
						})

						Convey("Should successfully create object using a mapping", func() {
							ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
							objs := util.MustStringify([]map[string]interface{}{
								{"user_id": account["id"].(string), "first_name": "John", "last_name": "Doe"},
							})
							resp, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs, Mapping: mappingName})
							So(err, ShouldBeNil)
							var m []map[string]interface{}
							err = util.FromJSON(resp.Objects, &m)
							So(err, ShouldBeNil)
							So(m, ShouldHaveLength, 1)
							So(m[0], ShouldContainKey, "user_id")
							So(m[0], ShouldContainKey, "first_name")
							So(m[0], ShouldContainKey, "last_name")
							So(m[0]["user_id"], ShouldEqual, account["id"])
							So(m[0]["first_name"], ShouldEqual, "John")
							So(m[0]["last_name"], ShouldEqual, "Doe")
						})
					})

					Convey("With local session", func() {
						Convey("Should return error if session does not exists locally and in session registry", func() {
							ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
							ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", "invalid"))
							objs := util.MustStringify([]map[string]interface{}{
								{"owner_id": account["id"].(string), "key": "mykey", "value": "myval"},
							})
							_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs})
							So(err, ShouldNotBeNil)
							m, err := util.JSONToMap(err.Error())
							So(err, ShouldBeNil)
							errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
							So(errs, ShouldHaveLength, 1)
							So(errs[0], ShouldResemble, map[string]interface{}{
								"status":  "404",
								"message": "session not found",
							})
						})

						Convey("Should return error when using a local, unregistered session not belonging to the authenticated account/caller", func() {
							sessionID := util.UUID4()
							rpc.dbSession.CreateUnregisteredSession(sessionID, account2["id"].(string))
							ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
							ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID))
							objs := util.MustStringify([]map[string]interface{}{
								{"owner_id": account["id"].(string), "key": "mykey", "value": "myval"},
							})
							_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs})
							So(err, ShouldNotBeNil)
							m, err := util.JSONToMap(err.Error())
							So(err, ShouldBeNil)
							errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
							So(errs, ShouldHaveLength, 1)
							So(errs[0], ShouldResemble, map[string]interface{}{
								"status":  "401",
								"message": "permission denied: you don't have permission to perform this operation",
							})

							rpc.dbSession.End(sessionID)
						})

						Convey("Should successfully create object using a local, unregistered session", func() {
							sessionID := util.UUID4()
							rpc.dbSession.CreateUnregisteredSession(sessionID, account["id"].(string))
							So(err, ShouldBeNil)
							ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
							ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID))
							objs := util.MustStringify([]map[string]interface{}{
								{"owner_id": account["id"].(string), "key": "mykey", "value": "myval"},
							})
							resp, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs})
							So(err, ShouldBeNil)
							var m []map[string]interface{}
							err = util.FromJSON(resp.Objects, &m)
							So(err, ShouldBeNil)
							So(m, ShouldHaveLength, 1)
							So(m[0]["bucket"], ShouldEqual, b.Name)
							So(m[0]["creator_id"], ShouldEqual, account["id"])
							So(m[0]["owner_id"], ShouldEqual, account["id"])
							rpc.dbSession.CommitEnd(sessionID)
						})
					})

					Convey("With remote session", func() {
						Convey("Should return error if authenticated user is not the owner of the session", func() {
							sessionID, err := rpc2.dbSession.CreateSession(util.UUID4(), "some_id")
							So(err, ShouldBeNil)

							token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
								"id":   account["client_id"],
								"type": oauth.TokenTypeApp,
								"iat":  time.Now().Unix(),
							})

							ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
							ctx2 := metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID, "authorization", "bearer "+token))
							_, err = rpc.CreateObjects(ctx2, &proto_rpc.CreateObjectsMsg{
								Bucket: b.Name,
							})
							So(err, ShouldNotBeNil)
							m, err := util.JSONToMap(err.Error())
							So(err, ShouldBeNil)
							errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
							So(errs, ShouldHaveLength, 1)
							So(errs[0], ShouldResemble, map[string]interface{}{
								"status":  "401",
								"message": "permission denied: you don't have permission to perform this operation",
							})
						})

						Convey("Should return error if session is not found in remote node", func() {
							sessionID, err := rpc2.dbSession.CreateSession(util.UUID4(), account["id"].(string))
							So(err, ShouldBeNil)

							// remove session manually from remote node
							rpc2.dbSession.RemoveAgent(sessionID)

							token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
								"id":   account["client_id"],
								"type": oauth.TokenTypeApp,
								"iat":  time.Now().Unix(),
							})

							ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
							ctx2 := metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID, "authorization", "bearer "+token))
							_, err = rpc.CreateObjects(ctx2, &proto_rpc.CreateObjectsMsg{
								Bucket: b.Name,
							})
							So(err, ShouldNotBeNil)
							m, err := util.JSONToMap(err.Error())
							So(err, ShouldBeNil)
							errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
							So(errs, ShouldHaveLength, 1)
							So(errs[0], ShouldResemble, map[string]interface{}{
								"status":  "404",
								"message": "session not found",
							})
						})
					})
				})
			})

			Convey(".GetObjects", func() {

				ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
				objs := util.MustStringify([]map[string]interface{}{
					{"owner_id": account["id"].(string), "key": "mykey", "value": "myval"},
					{"owner_id": account["id"].(string), "key": "mykey2", "value": "myval2"},
				})
				_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs})
				So(err, ShouldBeNil)

				Convey("Should return error if bucket name is not provided", func() {
					ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
					_, err := rpc.GetObjects(ctx, &proto_rpc.GetObjectMsg{})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "400",
						"field":   "bucket",
						"message": "bucket name is required",
					})
				})

				Convey("Should return error if bucket is not provided", func() {
					_, err := rpc.GetObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: ""})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "400",
						"field":   "bucket",
						"message": "bucket name is required",
					})
				})

				Convey("Should return error if bucket does not exists", func() {
					_, err := rpc.GetObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: "unknown"})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"field":   "bucket",
						"message": "bucket not found",
						"status":  "404",
					})
				})

				Convey("Should return error if owner does not exist", func() {
					_, err := rpc.GetObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: b.Name, Owner: "unknown"})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"message": "owner not found",
						"status":  "404",
						"field":   "owner",
					})
				})

				Convey("Should return error if creator does not exist", func() {
					_, err := rpc.GetObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: b.Name, Creator: "unknown"})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"message": "creator not found",
						"status":  "404",
						"field":   "creator",
					})
				})

				Convey("Should return error if authenticated user is not the owner of the queried object", func() {
					_, err := rpc.GetObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: b.Name, Owner: account2["id"].(string)})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "401",
						"message": "permission denied: you are not authorized to access objects belonging to the owner",
					})
				})

				Convey("Should return error if query parsing error occurred", func() {
					q := []byte(`{"unknown_field": "value"}`)
					_, err := rpc.GetObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: b.Name, Query: q, Owner: account["id"].(string)})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"message": "unknown query field: unknown_field",
						"status":  "400",
						"code":    "invalid_parameter",
						"field":   "query",
					})
				})

				Convey("Should successfully fetch objects", func() {
					resp, err := rpc.GetObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: b.Name, Owner: account["id"].(string)})
					So(err, ShouldBeNil)
					var m []map[string]interface{}
					err = util.FromJSON(resp.Objects, &m)
					So(err, ShouldBeNil)
					So(m, ShouldHaveLength, 2)
				})

				Convey("With mapping", func() {
					ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
					mappingName := util.RandString(5)
					mapping := `{ "user_id": "owner_id", "first_name": "key", "last_name": "value" }`
					_, err := rpc.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{Name: mappingName, Bucket: b.Name, Mapping: []byte(mapping)})
					So(err, ShouldBeNil)

					Convey("Should return error if mapping does not exist", func() {
						_, err := rpc.GetObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: b.Name, Owner: account["id"].(string), Mapping: "unknown"})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"status":  "404",
							"code":    "invalid_parameter",
							"field":   "mapping",
							"message": "mapping not found",
						})
					})

					Convey("Should successfully fetch objects with a mapping applied", func() {
						resp, err := rpc.GetObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: b.Name, Owner: account["id"].(string), Mapping: mappingName})
						So(err, ShouldBeNil)
						var m []map[string]interface{}
						err = util.FromJSON(resp.Objects, &m)
						So(err, ShouldBeNil)
						So(m, ShouldHaveLength, 2)
						So(m[0], ShouldContainKey, "user_id")
						So(m[0], ShouldContainKey, "first_name")
						So(m[0], ShouldContainKey, "last_name")
						So(m[1], ShouldContainKey, "user_id")
						So(m[1], ShouldContainKey, "first_name")
						So(m[1], ShouldContainKey, "last_name")
					})

				})

				Convey("With remote session", func() {

					Convey("Should return error if session is not found on local or remote registry", func() {
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						ctx2 := metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", "unknown_session"))
						_, err := rpc.GetObjects(ctx2, &proto_rpc.GetObjectMsg{Bucket: b.Name})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"message": "session not found",
							"status":  "404",
						})
					})

					Convey("Should return error if authenticated account/caller is not the owner of the session", func() {
						sessionID, err := rpc2.dbSession.CreateSession(util.UUID4(), "some_id")
						So(err, ShouldBeNil)

						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						ctx2 := metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID))
						_, err = rpc.GetObjects(ctx2, &proto_rpc.GetObjectMsg{Bucket: b.Name})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"message": "permission denied: you don't have permission to perform this operation",
							"status":  "401",
						})
					})

					Convey("Should return error if session exists in registry but not found in the remote node", func() {
						sessionID, err := rpc2.dbSession.CreateSession(util.UUID4(), account["id"].(string))
						So(err, ShouldBeNil)

						rpc2.dbSession.RemoveAgent(sessionID)

						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"id":   account["client_id"],
							"type": oauth.TokenTypeApp,
							"iat":  time.Now().Unix(),
						})

						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						ctx2 := metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID, "authorization", "bearer "+token))
						_, err = rpc.GetObjects(ctx2, &proto_rpc.GetObjectMsg{Bucket: b.Name})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"status":  "404",
							"message": "session not found",
						})
					})

					Convey("Should successfully get objects", func() {
						sessionID, err := rpc2.dbSession.CreateSession(util.UUID4(), account["id"].(string))
						So(err, ShouldBeNil)

						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"id":   account["client_id"],
							"type": oauth.TokenTypeApp,
							"iat":  time.Now().Unix(),
						})

						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						ctx2 := metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID, "authorization", "bearer "+token))
						resp, err := rpc.GetObjects(ctx2, &proto_rpc.GetObjectMsg{Bucket: b.Name})
						So(err, ShouldBeNil)
						var m []map[string]interface{}
						err = util.FromJSON(resp.Objects, &m)
						So(err, ShouldBeNil)
						So(m, ShouldHaveLength, 2)
						rpc2.dbSession.CommitEnd(sessionID)
					})
				})
			})

			Convey(".UpdateObjects", func() {

				ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
				imBucket := &proto_rpc.CreateBucketMsg{Name: util.RandString(5), Immutable: true}
				bucket, err := rpc.CreateBucket(ctx, imBucket)
				So(err, ShouldBeNil)
				So(bucket.Name, ShouldEqual, imBucket.Name)
				So(bucket.ID, ShouldHaveLength, 36)

				objs := []map[string]interface{}{
					{"owner_id": account["id"].(string), "key": util.RandString(5), "value": "myval"},
					{"owner_id": account["id"].(string), "key": util.RandString(5), "value": "myval2"},
				}
				_, err = rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: util.MustStringify(objs)})
				So(err, ShouldBeNil)

				Convey("Should return error if bucket name is provided", func() {
					ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
					_, err := rpc.UpdateObjects(ctx, &proto_rpc.UpdateObjectsMsg{})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "400",
						"field":   "bucket",
						"message": "bucket name is required",
					})
				})

				Convey("Should return error if bucket does not exist", func() {
					ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
					_, err := rpc.UpdateObjects(ctx, &proto_rpc.UpdateObjectsMsg{Bucket: "unknown"})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "404",
						"field":   "bucket",
						"message": "bucket not found",
					})
				})

				Convey("Should return error if bucket is immutable", func() {
					ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
					_, err := rpc.UpdateObjects(ctx, &proto_rpc.UpdateObjectsMsg{Bucket: imBucket.Name})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "400",
						"field":   "bucket",
						"message": "bucket is not mutable",
					})
				})

				Convey("Should return error if owner is not found", func() {
					ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
					_, err := rpc.UpdateObjects(ctx, &proto_rpc.UpdateObjectsMsg{Bucket: b.Name, Owner: "unknown"})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "404",
						"field":   "owner",
						"message": "owner not found",
					})
				})

				Convey("Should return error if creator is not found", func() {
					ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
					_, err := rpc.UpdateObjects(ctx, &proto_rpc.UpdateObjectsMsg{Bucket: b.Name, Creator: "unknown"})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "404",
						"field":   "creator",
						"message": "creator not found",
					})
				})

				Convey("Should return error if authenticated user is not authorized to access objects belonging to the owner", func() {
					ctx := context.WithValue(context.Background(), CtxAccount, account2["id"])
					_, err := rpc.UpdateObjects(ctx, &proto_rpc.UpdateObjectsMsg{Bucket: b.Name, Owner: account["id"].(string)})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "401",
						"message": "permission denied: you are not authorized to access objects belonging to the owner",
					})
				})

				Convey("Without mapping", func() {

					Convey("Should return parser error if query is invalid", func() {
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						_, err := rpc.UpdateObjects(ctx, &proto_rpc.UpdateObjectsMsg{
							Bucket: b.Name,
							Owner:  account["id"].(string),
							Query:  util.MustStringify(map[string]interface{}{"unknown_field": ""}),
							Update: util.MustStringify(map[string]interface{}{"value": "new value"}),
						})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"message": "unknown query field: unknown_field",
							"status":  "400",
							"code":    "invalid_parameter",
							"field":   "query",
						})
					})

					Convey("Should successfully update objects", func() {
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						resp, err := rpc.UpdateObjects(ctx, &proto_rpc.UpdateObjectsMsg{
							Bucket: b.Name,
							Owner:  account["id"].(string),
							Query:  util.MustStringify(map[string]interface{}{"key": objs[0]["key"]}),
							Update: util.MustStringify(map[string]interface{}{"value": "new value"}),
						})
						So(err, ShouldBeNil)
						So(resp.Affected, ShouldEqual, 1)

						var obj db.Object
						err = rpc.db.Where(db.Object{Key: objs[0]["key"].(string)}).Find(&obj).Error
						So(err, ShouldBeNil)
						So(obj.Value, ShouldEqual, "new value")
					})
				})

				Convey("With mapping", func() {
					Convey("Should return error if mapping does not exist", func() {
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						_, err := rpc.UpdateObjects(ctx, &proto_rpc.UpdateObjectsMsg{
							Bucket:  b.Name,
							Owner:   account["id"].(string),
							Mapping: "unknown",
						})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"field":   "mapping",
							"message": "mapping not found",
							"status":  "404",
							"code":    "invalid_parameter",
						})
					})
				})

				Convey("With session", func() {

					Convey("Should return error if session does not exist", func() {
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", "unknown"))
						_, err := rpc.UpdateObjects(ctx, &proto_rpc.UpdateObjectsMsg{
							Bucket: b.Name,
							Owner:  account["id"].(string),
							Query:  util.MustStringify(map[string]interface{}{"key": objs[0]["key"]}),
							Update: util.MustStringify(map[string]interface{}{"value": "new value"}),
						})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"status":  "404",
							"message": "session not found",
						})
					})

					Convey("Should successfully update object", func() {
						sessionID := rpc.dbSession.CreateUnregisteredSession("", account["id"].(string))
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID))
						resp, err := rpc.UpdateObjects(ctx, &proto_rpc.UpdateObjectsMsg{
							Bucket: b.Name,
							Owner:  account["id"].(string),
							Query:  util.MustStringify(map[string]interface{}{"key": objs[0]["key"]}),
							Update: util.MustStringify(map[string]interface{}{"value": "new value"}),
						})
						So(err, ShouldBeNil)
						So(resp.Affected, ShouldEqual, 1)
						rpc.dbSession.CommitEnd(sessionID)

						var obj db.Object
						err = rpc.db.Where(db.Object{Key: objs[0]["key"].(string)}).Find(&obj).Error
						So(err, ShouldBeNil)
						So(obj.Value, ShouldEqual, "new value")
					})

				})

				Convey("With remote session", func() {

					Convey("Should return error if remote session exist in registory but not in the remote node", func() {

						sessionID, err := rpc2.dbSession.CreateSession(util.UUID4(), account["id"].(string))
						So(err, ShouldBeNil)

						// remove session manually from remote node
						rpc2.dbSession.RemoveAgent(sessionID)

						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"id":   account["client_id"],
							"type": oauth.TokenTypeApp,
							"iat":  time.Now().Unix(),
						})

						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						ctx2 := metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID, "authorization", "bearer "+token))
						_, err = rpc.UpdateObjects(ctx2, &proto_rpc.UpdateObjectsMsg{
							Bucket: b.Name,
							Owner:  account["id"].(string),
							Query:  util.MustStringify(map[string]interface{}{"key": objs[0]["key"]}),
							Update: util.MustStringify(map[string]interface{}{"value": "new value"}),
						})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"status":  "404",
							"message": "session not found",
						})
					})

					Convey("Should return error if session is not owned by authenticated account", func() {
						sessionID, err := rpc2.dbSession.CreateSession(util.UUID4(), account2["id"].(string))
						So(err, ShouldBeNil)
						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"id":   account["client_id"],
							"type": oauth.TokenTypeApp,
							"iat":  time.Now().Unix(),
						})

						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						ctx2 := metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID, "authorization", "bearer "+token))
						_, err = rpc.UpdateObjects(ctx2, &proto_rpc.UpdateObjectsMsg{
							Bucket: b.Name,
							Owner:  account["id"].(string),
							Query:  util.MustStringify(map[string]interface{}{"key": objs[0]["key"]}),
							Update: util.MustStringify(map[string]interface{}{"value": "new value"}),
						})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"status":  "401",
							"message": "permission denied: you don't have permission to perform this operation",
						})
						rpc2.dbSession.RollbackEnd(sessionID)
					})
				})
			})

			Convey(".CountObjects", func() {
				ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
				objs := util.MustStringify([]map[string]interface{}{
					{"owner_id": account["id"].(string), "key": "mykey", "value": "myval"},
					{"owner_id": account["id"].(string), "key": "mykey2", "value": "myval2"},
				})
				_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: objs})
				So(err, ShouldBeNil)

				Convey("Should return error if bucket is not provided", func() {
					_, err := rpc.CountObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: ""})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "400",
						"field":   "bucket",
						"message": "bucket name is required",
					})
				})

				Convey("Should return error if authenticated user is not authorized to access objects belonging to the owner", func() {
					_, err := rpc.CountObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: b.Name, Owner: account2["id"].(string)})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "401",
						"message": "permission denied: you are not authorized to access objects belonging to the owner",
					})
				})

				Convey("Should return error if bucket does not exists", func() {
					_, err := rpc.CountObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: "unknown"})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"field":   "bucket",
						"message": "bucket not found",
						"status":  "404",
					})
				})

				Convey("Should return error if owner does not exist", func() {
					_, err := rpc.CountObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: b.Name, Owner: "unknown"})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"message": "owner not found",
						"status":  "404",
						"field":   "owner",
					})
				})

				Convey("Should return error if creator does not exist", func() {
					_, err := rpc.GetObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: b.Name, Creator: "unknown"})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"message": "creator not found",
						"status":  "404",
						"field":   "creator",
					})
				})

				Convey("Should successfully return expected count", func() {
					resp, err := rpc.CountObjects(ctx, &proto_rpc.GetObjectMsg{Bucket: b.Name, Owner: account["id"].(string)})
					So(err, ShouldBeNil)
					So(resp.Count, ShouldEqual, 2)
				})

				Convey("With remote session", func() {
					Convey("Should return error if session exists in registry but not found in the remote node", func() {
						sessionID, err := rpc2.dbSession.CreateSession(util.UUID4(), account["id"].(string))
						So(err, ShouldBeNil)

						rpc2.dbSession.RemoveAgent(sessionID)

						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"id":   account["client_id"],
							"type": oauth.TokenTypeApp,
							"iat":  time.Now().Unix(),
						})

						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						ctx2 := metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID, "authorization", "bearer "+token))
						_, err = rpc.CountObjects(ctx2, &proto_rpc.GetObjectMsg{Bucket: b.Name})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"status":  "404",
							"message": "session not found",
						})
					})

					Convey("Should successfully get object count", func() {
						sessionID, err := rpc2.dbSession.CreateSession(util.UUID4(), account["id"].(string))
						So(err, ShouldBeNil)

						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"id":   account["client_id"],
							"type": oauth.TokenTypeApp,
							"iat":  time.Now().Unix(),
						})

						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						ctx2 := metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID, "authorization", "bearer "+token))
						resp, err := rpc.CountObjects(ctx2, &proto_rpc.GetObjectMsg{Bucket: b.Name})
						So(err, ShouldBeNil)
						So(resp.Count, ShouldEqual, 2)
						rpc2.dbSession.CommitEnd(sessionID)
					})
				})
			})

			Convey(".DeleteObjects", func() {

				ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
				objs := []map[string]interface{}{
					{"owner_id": account["id"].(string), "key": "obj_to_del", "value": "myval2"},
				}
				_, err := rpc.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Bucket: b.Name, Objects: util.MustStringify(objs)})
				So(err, ShouldBeNil)

				Convey("Should return error if bucket is not provided", func() {
					_, err := rpc.DeleteObjects(ctx, &proto_rpc.DeleteObjectsMsg{Bucket: ""})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "400",
						"field":   "bucket",
						"message": "bucket name is required",
					})
				})

				Convey("Should return error if bucket does not exists", func() {
					_, err := rpc.DeleteObjects(ctx, &proto_rpc.DeleteObjectsMsg{Bucket: "unknown"})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"field":   "bucket",
						"message": "bucket not found",
						"status":  "404",
					})
				})

				Convey("Should return error if bucket is not mutable", func() {
					_, err := rpc.DeleteObjects(ctx, &proto_rpc.DeleteObjectsMsg{Bucket: b2.Name})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "400",
						"field":   "bucket",
						"message": "bucket is not mutable",
					})
				})

				Convey("Should return error if owner does not exist", func() {
					_, err := rpc.DeleteObjects(ctx, &proto_rpc.DeleteObjectsMsg{Bucket: b.Name, Owner: "unknown"})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "404",
						"field":   "owner",
						"message": "owner not found",
					})
				})

				Convey("Should return error if authenticated user is not authorized to access objects belonging to the owner", func() {
					_, err := rpc.DeleteObjects(ctx, &proto_rpc.DeleteObjectsMsg{Bucket: b.Name, Owner: account2["id"].(string)})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"status":  "401",
						"message": "permission denied: you are not authorized to access objects belonging to the owner",
					})
				})

				Convey("Should return parser error if query is invalid", func() {
					_, err := rpc.DeleteObjects(ctx, &proto_rpc.DeleteObjectsMsg{Bucket: b.Name, Query: []byte(`{"unknown_key": "value"}`), Owner: account["id"].(string)})
					So(err, ShouldNotBeNil)
					m, err := util.JSONToMap(err.Error())
					So(err, ShouldBeNil)
					errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
					So(errs, ShouldHaveLength, 1)
					So(errs[0], ShouldResemble, map[string]interface{}{
						"code":    "invalid_parameter",
						"field":   "query",
						"message": "unknown query field: unknown_key",
						"status":  "400",
					})
				})

				Convey("Should successfully delete object", func() {
					resp, err := rpc.DeleteObjects(ctx, &proto_rpc.DeleteObjectsMsg{
						Bucket: b.Name,
						Owner:  account["id"].(string),
						Query: util.MustStringify(map[string]interface{}{
							"key": "obj_to_del",
						}),
					})
					So(err, ShouldBeNil)
					So(resp.Affected, ShouldEqual, 1)
				})

				Convey("With mapping", func() {

					Convey("Should return error if mapping does not exist", func() {
						_, err := rpc.DeleteObjects(ctx, &proto_rpc.DeleteObjectsMsg{Bucket: b.Name, Mapping: "unknown", Owner: account["id"].(string)})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"status":  "404",
							"code":    "invalid_parameter",
							"field":   "mapping",
							"message": "mapping not found",
						})
					})

					Convey("Should successfully delete object", func() {
						resp, err := rpc.DeleteObjects(ctx, &proto_rpc.DeleteObjectsMsg{
							Bucket:  b.Name,
							Owner:   account["id"].(string),
							Mapping: mapping.Name,
							Query: util.MustStringify(map[string]interface{}{
								"name": "obj_to_del",
							}),
						})
						So(err, ShouldBeNil)
						So(resp.Affected, ShouldEqual, 1)
					})
				})

				Convey("With remote session", func() {
					Convey("Should return error if session exists in registry but not found in the remote node", func() {
						sessionID, err := rpc2.dbSession.CreateSession(util.UUID4(), account["id"].(string))
						So(err, ShouldBeNil)

						rpc2.dbSession.RemoveAgent(sessionID)

						token, _ := oauth.MakeToken(oauth.SigningSecret, map[string]interface{}{
							"id":   account["client_id"],
							"type": oauth.TokenTypeApp,
							"iat":  time.Now().Unix(),
						})

						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						ctx2 := metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID, "authorization", "bearer "+token))
						_, err = rpc.DeleteObjects(ctx2, &proto_rpc.DeleteObjectsMsg{Bucket: b.Name})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"status":  "404",
							"message": "session not found",
						})
					})
				})

				Convey("With local session", func() {
					Convey("Should return error if session does not exists locally or in registry", func() {
						ctx := context.WithValue(context.Background(), CtxAccount, account["id"])
						ctx2 := metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", "unknown"))
						_, err = rpc.DeleteObjects(ctx2, &proto_rpc.DeleteObjectsMsg{Bucket: b.Name})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"status":  "404",
							"message": "session not found",
						})
					})

					Convey("Should return permission error if session is not owned by authenticated account", func() {
						sessionID := rpc.dbSession.CreateUnregisteredSession("", account["id"].(string))
						ctx := context.WithValue(context.Background(), CtxAccount, account2["id"])
						ctx2 := metadata.NewIncomingContext(ctx, metadata.Pairs("session_id", sessionID))
						_, err = rpc.DeleteObjects(ctx2, &proto_rpc.DeleteObjectsMsg{Bucket: b.Name})
						So(err, ShouldNotBeNil)
						m, err := util.JSONToMap(err.Error())
						So(err, ShouldBeNil)
						errs := m["Errors"].(map[string]interface{})["errors"].([]interface{})
						So(errs, ShouldHaveLength, 1)
						So(errs[0], ShouldResemble, map[string]interface{}{
							"message": "permission denied: you don't have permission to perform this operation",
							"status":  "401",
						})
					})
				})
			})
		})
	})
}

func TestOrderByString(t *testing.T) {
	Convey(".orderByToString", t, func() {
		Convey("Should expected result", func() {
			str := orderByToString([]*proto_rpc.OrderBy{
				{Field: "field_a", Order: 0},
				{Field: "field_b", Order: 1},
			})
			So(str, ShouldEqual, "field_a desc, field_b asc")
		})
	})
}
