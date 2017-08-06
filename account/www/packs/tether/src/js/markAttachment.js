/* globals Tether */
Tether.modules.push({
    initialize: function () {
        var _this = this;
        this.markers = {};
        ['target', 'element'].forEach(function (type) {
            var el = document.createElement('div');
            el.className = _this.getClass(type + "-marker");
            var dot = document.createElement('div');
            dot.className = _this.getClass('marker-dot');
            el.appendChild(dot);
            _this[type].appendChild(el);
            _this.markers[type] = { dot: dot, el: el };
        });
    },
    position: function (_b) {
        var manualOffset = _b.manualOffset, manualTargetOffset = _b.manualTargetOffset;
        var offsets = {
            element: manualOffset,
            target: manualTargetOffset
        };
        for (var type in offsets) {
            var offset = offsets[type];
            for (var side in offset) {
                var val = offset[side];
                var notString = typeof val !== 'string';
                if (notString ||
                    val.indexOf('%') === -1 &&
                        val.indexOf('px') === -1) {
                    val += 'px';
                }
                if (this.markers[type].dot.style[side] !== val) {
                    this.markers[type].dot.style[side] = val;
                }
            }
        }
        return true;
    }
});
//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbInd3dy9wYWNrcy90ZXRoZXIvc3JjL2pzL21hcmtBdHRhY2htZW50LmpzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUFBLG9CQUFvQjtBQUVwQixNQUFNLENBQUMsT0FBTyxDQUFDLElBQUksQ0FBQztJQUNsQixVQUFVO1FBQVYsaUJBZUM7UUFkQyxJQUFJLENBQUMsT0FBTyxHQUFHLEVBQUUsQ0FBQztRQUVsQixDQUFDLFFBQVEsRUFBRSxTQUFTLENBQUMsQ0FBQyxPQUFPLENBQUMsVUFBQSxJQUFJO1lBQ2hDLElBQU0sRUFBRSxHQUFHLFFBQVEsQ0FBQyxhQUFhLENBQUMsS0FBSyxDQUFDLENBQUM7WUFDekMsRUFBRSxDQUFDLFNBQVMsR0FBRyxLQUFJLENBQUMsUUFBUSxDQUFLLElBQUksWUFBVSxDQUFDLENBQUM7WUFFakQsSUFBTSxHQUFHLEdBQUcsUUFBUSxDQUFDLGFBQWEsQ0FBQyxLQUFLLENBQUMsQ0FBQztZQUMxQyxHQUFHLENBQUMsU0FBUyxHQUFHLEtBQUksQ0FBQyxRQUFRLENBQUMsWUFBWSxDQUFDLENBQUM7WUFDNUMsRUFBRSxDQUFDLFdBQVcsQ0FBQyxHQUFHLENBQUMsQ0FBQztZQUVwQixLQUFJLENBQUMsSUFBSSxDQUFDLENBQUMsV0FBVyxDQUFDLEVBQUUsQ0FBQyxDQUFDO1lBRTNCLEtBQUksQ0FBQyxPQUFPLENBQUMsSUFBSSxDQUFDLEdBQUcsRUFBQyxHQUFHLEtBQUEsRUFBRSxFQUFFLElBQUEsRUFBQyxDQUFDO1FBQ2pDLENBQUMsQ0FBQyxDQUFDO0lBQ0wsQ0FBQztJQUVELFFBQVEsWUFBQyxFQUFrQztZQUFqQyw4QkFBWSxFQUFFLDBDQUFrQjtRQUN4QyxJQUFNLE9BQU8sR0FBRztZQUNkLE9BQU8sRUFBRSxZQUFZO1lBQ3JCLE1BQU0sRUFBRSxrQkFBa0I7U0FDM0IsQ0FBQztRQUVGLEdBQUcsQ0FBQyxDQUFDLElBQUksSUFBSSxJQUFJLE9BQU8sQ0FBQyxDQUFDLENBQUM7WUFDekIsSUFBTSxNQUFNLEdBQUcsT0FBTyxDQUFDLElBQUksQ0FBQyxDQUFDO1lBQzdCLEdBQUcsQ0FBQyxDQUFDLElBQUksSUFBSSxJQUFJLE1BQU0sQ0FBQyxDQUFDLENBQUM7Z0JBQ3hCLElBQUksR0FBRyxHQUFHLE1BQU0sQ0FBQyxJQUFJLENBQUMsQ0FBQztnQkFDdkIsSUFBTSxTQUFTLEdBQUcsT0FBTyxHQUFHLEtBQUssUUFBUSxDQUFDO2dCQUMxQyxFQUFFLENBQUMsQ0FBQyxTQUFTO29CQUNULEdBQUcsQ0FBQyxPQUFPLENBQUMsR0FBRyxDQUFDLEtBQUssQ0FBQyxDQUFDO3dCQUN2QixHQUFHLENBQUMsT0FBTyxDQUFDLElBQUksQ0FBQyxLQUFLLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQztvQkFDN0IsR0FBRyxJQUFJLElBQUksQ0FBQztnQkFDZCxDQUFDO2dCQUVELEVBQUUsQ0FBQyxDQUFDLElBQUksQ0FBQyxPQUFPLENBQUMsSUFBSSxDQUFDLENBQUMsR0FBRyxDQUFDLEtBQUssQ0FBQyxJQUFJLENBQUMsS0FBSyxHQUFHLENBQUMsQ0FBQyxDQUFDO29CQUMvQyxJQUFJLENBQUMsT0FBTyxDQUFDLElBQUksQ0FBQyxDQUFDLEdBQUcsQ0FBQyxLQUFLLENBQUMsSUFBSSxDQUFDLEdBQUcsR0FBRyxDQUFDO2dCQUMzQyxDQUFDO1lBQ0gsQ0FBQztRQUNILENBQUM7UUFFRCxNQUFNLENBQUMsSUFBSSxDQUFDO0lBQ2QsQ0FBQztDQUNGLENBQUMsQ0FBQyIsImZpbGUiOiJ3d3cvcGFja3MvdGV0aGVyL3NyYy9qcy9tYXJrQXR0YWNobWVudC5qcyIsInNvdXJjZXNDb250ZW50IjpbIi8qIGdsb2JhbHMgVGV0aGVyICovXG5cblRldGhlci5tb2R1bGVzLnB1c2goe1xuICBpbml0aWFsaXplKCkge1xuICAgIHRoaXMubWFya2VycyA9IHt9O1xuXG4gICAgWyd0YXJnZXQnLCAnZWxlbWVudCddLmZvckVhY2godHlwZSA9PiB7XG4gICAgICBjb25zdCBlbCA9IGRvY3VtZW50LmNyZWF0ZUVsZW1lbnQoJ2RpdicpO1xuICAgICAgZWwuY2xhc3NOYW1lID0gdGhpcy5nZXRDbGFzcyhgJHsgdHlwZSB9LW1hcmtlcmApO1xuXG4gICAgICBjb25zdCBkb3QgPSBkb2N1bWVudC5jcmVhdGVFbGVtZW50KCdkaXYnKTtcbiAgICAgIGRvdC5jbGFzc05hbWUgPSB0aGlzLmdldENsYXNzKCdtYXJrZXItZG90Jyk7XG4gICAgICBlbC5hcHBlbmRDaGlsZChkb3QpO1xuXG4gICAgICB0aGlzW3R5cGVdLmFwcGVuZENoaWxkKGVsKTtcblxuICAgICAgdGhpcy5tYXJrZXJzW3R5cGVdID0ge2RvdCwgZWx9O1xuICAgIH0pO1xuICB9LFxuXG4gIHBvc2l0aW9uKHttYW51YWxPZmZzZXQsIG1hbnVhbFRhcmdldE9mZnNldH0pIHtcbiAgICBjb25zdCBvZmZzZXRzID0ge1xuICAgICAgZWxlbWVudDogbWFudWFsT2Zmc2V0LFxuICAgICAgdGFyZ2V0OiBtYW51YWxUYXJnZXRPZmZzZXRcbiAgICB9O1xuXG4gICAgZm9yIChsZXQgdHlwZSBpbiBvZmZzZXRzKSB7XG4gICAgICBjb25zdCBvZmZzZXQgPSBvZmZzZXRzW3R5cGVdO1xuICAgICAgZm9yIChsZXQgc2lkZSBpbiBvZmZzZXQpIHtcbiAgICAgICAgbGV0IHZhbCA9IG9mZnNldFtzaWRlXTtcbiAgICAgICAgY29uc3Qgbm90U3RyaW5nID0gdHlwZW9mIHZhbCAhPT0gJ3N0cmluZyc7XG4gICAgICAgIGlmIChub3RTdHJpbmcgfHxcbiAgICAgICAgICAgIHZhbC5pbmRleE9mKCclJykgPT09IC0xICYmXG4gICAgICAgICAgICB2YWwuaW5kZXhPZigncHgnKSA9PT0gLTEpIHtcbiAgICAgICAgICB2YWwgKz0gJ3B4JztcbiAgICAgICAgfVxuXG4gICAgICAgIGlmICh0aGlzLm1hcmtlcnNbdHlwZV0uZG90LnN0eWxlW3NpZGVdICE9PSB2YWwpIHtcbiAgICAgICAgICB0aGlzLm1hcmtlcnNbdHlwZV0uZG90LnN0eWxlW3NpZGVdID0gdmFsO1xuICAgICAgICB9XG4gICAgICB9XG4gICAgfVxuXG4gICAgcmV0dXJuIHRydWU7XG4gIH1cbn0pO1xuIl19

//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbInd3dy9wYWNrcy90ZXRoZXIvc3JjL2pzL21hcmtBdHRhY2htZW50LmpzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUFBLG9CQUFvQjtBQUNwQixNQUFNLENBQUMsT0FBTyxDQUFDLElBQUksQ0FBQztJQUNoQixVQUFVLEVBQUU7UUFDUixJQUFJLEtBQUssR0FBRyxJQUFJLENBQUM7UUFDakIsSUFBSSxDQUFDLE9BQU8sR0FBRyxFQUFFLENBQUM7UUFDbEIsQ0FBQyxRQUFRLEVBQUUsU0FBUyxDQUFDLENBQUMsT0FBTyxDQUFDLFVBQVUsSUFBSTtZQUN4QyxJQUFJLEVBQUUsR0FBRyxRQUFRLENBQUMsYUFBYSxDQUFDLEtBQUssQ0FBQyxDQUFDO1lBQ3ZDLEVBQUUsQ0FBQyxTQUFTLEdBQUcsS0FBSyxDQUFDLFFBQVEsQ0FBQyxJQUFJLEdBQUcsU0FBUyxDQUFDLENBQUM7WUFDaEQsSUFBSSxHQUFHLEdBQUcsUUFBUSxDQUFDLGFBQWEsQ0FBQyxLQUFLLENBQUMsQ0FBQztZQUN4QyxHQUFHLENBQUMsU0FBUyxHQUFHLEtBQUssQ0FBQyxRQUFRLENBQUMsWUFBWSxDQUFDLENBQUM7WUFDN0MsRUFBRSxDQUFDLFdBQVcsQ0FBQyxHQUFHLENBQUMsQ0FBQztZQUNwQixLQUFLLENBQUMsSUFBSSxDQUFDLENBQUMsV0FBVyxDQUFDLEVBQUUsQ0FBQyxDQUFDO1lBQzVCLEtBQUssQ0FBQyxPQUFPLENBQUMsSUFBSSxDQUFDLEdBQUcsRUFBRSxHQUFHLEVBQUUsR0FBRyxFQUFFLEVBQUUsRUFBRSxFQUFFLEVBQUUsQ0FBQztRQUMvQyxDQUFDLENBQUMsQ0FBQztJQUNQLENBQUM7SUFDRCxRQUFRLEVBQUUsVUFBVSxFQUFFO1FBQ2xCLElBQUksWUFBWSxHQUFHLEVBQUUsQ0FBQyxZQUFZLEVBQUUsa0JBQWtCLEdBQUcsRUFBRSxDQUFDLGtCQUFrQixDQUFDO1FBQy9FLElBQUksT0FBTyxHQUFHO1lBQ1YsT0FBTyxFQUFFLFlBQVk7WUFDckIsTUFBTSxFQUFFLGtCQUFrQjtTQUM3QixDQUFDO1FBQ0YsR0FBRyxDQUFDLENBQUMsSUFBSSxJQUFJLElBQUksT0FBTyxDQUFDLENBQUMsQ0FBQztZQUN2QixJQUFJLE1BQU0sR0FBRyxPQUFPLENBQUMsSUFBSSxDQUFDLENBQUM7WUFDM0IsR0FBRyxDQUFDLENBQUMsSUFBSSxJQUFJLElBQUksTUFBTSxDQUFDLENBQUMsQ0FBQztnQkFDdEIsSUFBSSxHQUFHLEdBQUcsTUFBTSxDQUFDLElBQUksQ0FBQyxDQUFDO2dCQUN2QixJQUFJLFNBQVMsR0FBRyxPQUFPLEdBQUcsS0FBSyxRQUFRLENBQUM7Z0JBQ3hDLEVBQUUsQ0FBQyxDQUFDLFNBQVM7b0JBQ1QsR0FBRyxDQUFDLE9BQU8sQ0FBQyxHQUFHLENBQUMsS0FBSyxDQUFDLENBQUM7d0JBQ25CLEdBQUcsQ0FBQyxPQUFPLENBQUMsSUFBSSxDQUFDLEtBQUssQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDO29CQUMvQixHQUFHLElBQUksSUFBSSxDQUFDO2dCQUNoQixDQUFDO2dCQUNELEVBQUUsQ0FBQyxDQUFDLElBQUksQ0FBQyxPQUFPLENBQUMsSUFBSSxDQUFDLENBQUMsR0FBRyxDQUFDLEtBQUssQ0FBQyxJQUFJLENBQUMsS0FBSyxHQUFHLENBQUMsQ0FBQyxDQUFDO29CQUM3QyxJQUFJLENBQUMsT0FBTyxDQUFDLElBQUksQ0FBQyxDQUFDLEdBQUcsQ0FBQyxLQUFLLENBQUMsSUFBSSxDQUFDLEdBQUcsR0FBRyxDQUFDO2dCQUM3QyxDQUFDO1lBQ0wsQ0FBQztRQUNMLENBQUM7UUFDRCxNQUFNLENBQUMsSUFBSSxDQUFDO0lBQ2hCLENBQUM7Q0FDSixDQUFDLENBQUM7QUFFSCxtbkhBQW1uSCIsImZpbGUiOiJ3d3cvcGFja3MvdGV0aGVyL3NyYy9qcy9tYXJrQXR0YWNobWVudC5qcyIsInNvdXJjZXNDb250ZW50IjpbIi8qIGdsb2JhbHMgVGV0aGVyICovXG5UZXRoZXIubW9kdWxlcy5wdXNoKHtcbiAgICBpbml0aWFsaXplOiBmdW5jdGlvbiAoKSB7XG4gICAgICAgIHZhciBfdGhpcyA9IHRoaXM7XG4gICAgICAgIHRoaXMubWFya2VycyA9IHt9O1xuICAgICAgICBbJ3RhcmdldCcsICdlbGVtZW50J10uZm9yRWFjaChmdW5jdGlvbiAodHlwZSkge1xuICAgICAgICAgICAgdmFyIGVsID0gZG9jdW1lbnQuY3JlYXRlRWxlbWVudCgnZGl2Jyk7XG4gICAgICAgICAgICBlbC5jbGFzc05hbWUgPSBfdGhpcy5nZXRDbGFzcyh0eXBlICsgXCItbWFya2VyXCIpO1xuICAgICAgICAgICAgdmFyIGRvdCA9IGRvY3VtZW50LmNyZWF0ZUVsZW1lbnQoJ2RpdicpO1xuICAgICAgICAgICAgZG90LmNsYXNzTmFtZSA9IF90aGlzLmdldENsYXNzKCdtYXJrZXItZG90Jyk7XG4gICAgICAgICAgICBlbC5hcHBlbmRDaGlsZChkb3QpO1xuICAgICAgICAgICAgX3RoaXNbdHlwZV0uYXBwZW5kQ2hpbGQoZWwpO1xuICAgICAgICAgICAgX3RoaXMubWFya2Vyc1t0eXBlXSA9IHsgZG90OiBkb3QsIGVsOiBlbCB9O1xuICAgICAgICB9KTtcbiAgICB9LFxuICAgIHBvc2l0aW9uOiBmdW5jdGlvbiAoX2IpIHtcbiAgICAgICAgdmFyIG1hbnVhbE9mZnNldCA9IF9iLm1hbnVhbE9mZnNldCwgbWFudWFsVGFyZ2V0T2Zmc2V0ID0gX2IubWFudWFsVGFyZ2V0T2Zmc2V0O1xuICAgICAgICB2YXIgb2Zmc2V0cyA9IHtcbiAgICAgICAgICAgIGVsZW1lbnQ6IG1hbnVhbE9mZnNldCxcbiAgICAgICAgICAgIHRhcmdldDogbWFudWFsVGFyZ2V0T2Zmc2V0XG4gICAgICAgIH07XG4gICAgICAgIGZvciAodmFyIHR5cGUgaW4gb2Zmc2V0cykge1xuICAgICAgICAgICAgdmFyIG9mZnNldCA9IG9mZnNldHNbdHlwZV07XG4gICAgICAgICAgICBmb3IgKHZhciBzaWRlIGluIG9mZnNldCkge1xuICAgICAgICAgICAgICAgIHZhciB2YWwgPSBvZmZzZXRbc2lkZV07XG4gICAgICAgICAgICAgICAgdmFyIG5vdFN0cmluZyA9IHR5cGVvZiB2YWwgIT09ICdzdHJpbmcnO1xuICAgICAgICAgICAgICAgIGlmIChub3RTdHJpbmcgfHxcbiAgICAgICAgICAgICAgICAgICAgdmFsLmluZGV4T2YoJyUnKSA9PT0gLTEgJiZcbiAgICAgICAgICAgICAgICAgICAgICAgIHZhbC5pbmRleE9mKCdweCcpID09PSAtMSkge1xuICAgICAgICAgICAgICAgICAgICB2YWwgKz0gJ3B4JztcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgaWYgKHRoaXMubWFya2Vyc1t0eXBlXS5kb3Quc3R5bGVbc2lkZV0gIT09IHZhbCkge1xuICAgICAgICAgICAgICAgICAgICB0aGlzLm1hcmtlcnNbdHlwZV0uZG90LnN0eWxlW3NpZGVdID0gdmFsO1xuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuICAgICAgICByZXR1cm4gdHJ1ZTtcbiAgICB9XG59KTtcblxuLy8jIHNvdXJjZU1hcHBpbmdVUkw9ZGF0YTphcHBsaWNhdGlvbi9qc29uO2NoYXJzZXQ9dXRmODtiYXNlNjQsZXlKMlpYSnphVzl1SWpvekxDSnpiM1Z5WTJWeklqcGJJbmQzZHk5d1lXTnJjeTkwWlhSb1pYSXZjM0pqTDJwekwyMWhjbXRCZEhSaFkyaHRaVzUwTG1weklsMHNJbTVoYldWeklqcGJYU3dpYldGd2NHbHVaM01pT2lKQlFVRkJMRzlDUVVGdlFqdEJRVVZ3UWl4TlFVRk5MRU5CUVVNc1QwRkJUeXhEUVVGRExFbEJRVWtzUTBGQlF6dEpRVU5zUWl4VlFVRlZPMUZCUVZZc2FVSkJaVU03VVVGa1F5eEpRVUZKTEVOQlFVTXNUMEZCVHl4SFFVRkhMRVZCUVVVc1EwRkJRenRSUVVWc1FpeERRVUZETEZGQlFWRXNSVUZCUlN4VFFVRlRMRU5CUVVNc1EwRkJReXhQUVVGUExFTkJRVU1zVlVGQlFTeEpRVUZKTzFsQlEyaERMRWxCUVUwc1JVRkJSU3hIUVVGSExGRkJRVkVzUTBGQlF5eGhRVUZoTEVOQlFVTXNTMEZCU3l4RFFVRkRMRU5CUVVNN1dVRkRla01zUlVGQlJTeERRVUZETEZOQlFWTXNSMEZCUnl4TFFVRkpMRU5CUVVNc1VVRkJVU3hEUVVGTExFbEJRVWtzV1VGQlZTeERRVUZETEVOQlFVTTdXVUZGYWtRc1NVRkJUU3hIUVVGSExFZEJRVWNzVVVGQlVTeERRVUZETEdGQlFXRXNRMEZCUXl4TFFVRkxMRU5CUVVNc1EwRkJRenRaUVVNeFF5eEhRVUZITEVOQlFVTXNVMEZCVXl4SFFVRkhMRXRCUVVrc1EwRkJReXhSUVVGUkxFTkJRVU1zV1VGQldTeERRVUZETEVOQlFVTTdXVUZETlVNc1JVRkJSU3hEUVVGRExGZEJRVmNzUTBGQlF5eEhRVUZITEVOQlFVTXNRMEZCUXp0WlFVVndRaXhMUVVGSkxFTkJRVU1zU1VGQlNTeERRVUZETEVOQlFVTXNWMEZCVnl4RFFVRkRMRVZCUVVVc1EwRkJReXhEUVVGRE8xbEJSVE5DTEV0QlFVa3NRMEZCUXl4UFFVRlBMRU5CUVVNc1NVRkJTU3hEUVVGRExFZEJRVWNzUlVGQlF5eEhRVUZITEV0QlFVRXNSVUZCUlN4RlFVRkZMRWxCUVVFc1JVRkJReXhEUVVGRE8xRkJRMnBETEVOQlFVTXNRMEZCUXl4RFFVRkRPMGxCUTB3c1EwRkJRenRKUVVWRUxGRkJRVkVzV1VGQlF5eEZRVUZyUXp0WlFVRnFReXc0UWtGQldTeEZRVUZGTERCRFFVRnJRanRSUVVONFF5eEpRVUZOTEU5QlFVOHNSMEZCUnp0WlFVTmtMRTlCUVU4c1JVRkJSU3haUVVGWk8xbEJRM0pDTEUxQlFVMHNSVUZCUlN4clFrRkJhMEk3VTBGRE0wSXNRMEZCUXp0UlFVVkdMRWRCUVVjc1EwRkJReXhEUVVGRExFbEJRVWtzU1VGQlNTeEpRVUZKTEU5QlFVOHNRMEZCUXl4RFFVRkRMRU5CUVVNN1dVRkRla0lzU1VGQlRTeE5RVUZOTEVkQlFVY3NUMEZCVHl4RFFVRkRMRWxCUVVrc1EwRkJReXhEUVVGRE8xbEJRemRDTEVkQlFVY3NRMEZCUXl4RFFVRkRMRWxCUVVrc1NVRkJTU3hKUVVGSkxFMUJRVTBzUTBGQlF5eERRVUZETEVOQlFVTTdaMEpCUTNoQ0xFbEJRVWtzUjBGQlJ5eEhRVUZITEUxQlFVMHNRMEZCUXl4SlFVRkpMRU5CUVVNc1EwRkJRenRuUWtGRGRrSXNTVUZCVFN4VFFVRlRMRWRCUVVjc1QwRkJUeXhIUVVGSExFdEJRVXNzVVVGQlVTeERRVUZETzJkQ1FVTXhReXhGUVVGRkxFTkJRVU1zUTBGQlF5eFRRVUZUTzI5Q1FVTlVMRWRCUVVjc1EwRkJReXhQUVVGUExFTkJRVU1zUjBGQlJ5eERRVUZETEV0QlFVc3NRMEZCUXl4RFFVRkRPM2RDUVVOMlFpeEhRVUZITEVOQlFVTXNUMEZCVHl4RFFVRkRMRWxCUVVrc1EwRkJReXhMUVVGTExFTkJRVU1zUTBGQlF5eERRVUZETEVOQlFVTXNRMEZCUXp0dlFrRkROMElzUjBGQlJ5eEpRVUZKTEVsQlFVa3NRMEZCUXp0blFrRkRaQ3hEUVVGRE8yZENRVVZFTEVWQlFVVXNRMEZCUXl4RFFVRkRMRWxCUVVrc1EwRkJReXhQUVVGUExFTkJRVU1zU1VGQlNTeERRVUZETEVOQlFVTXNSMEZCUnl4RFFVRkRMRXRCUVVzc1EwRkJReXhKUVVGSkxFTkJRVU1zUzBGQlN5eEhRVUZITEVOQlFVTXNRMEZCUXl4RFFVRkRPMjlDUVVNdlF5eEpRVUZKTEVOQlFVTXNUMEZCVHl4RFFVRkRMRWxCUVVrc1EwRkJReXhEUVVGRExFZEJRVWNzUTBGQlF5eExRVUZMTEVOQlFVTXNTVUZCU1N4RFFVRkRMRWRCUVVjc1IwRkJSeXhEUVVGRE8yZENRVU16UXl4RFFVRkRPMWxCUTBnc1EwRkJRenRSUVVOSUxFTkJRVU03VVVGRlJDeE5RVUZOTEVOQlFVTXNTVUZCU1N4RFFVRkRPMGxCUTJRc1EwRkJRenREUVVOR0xFTkJRVU1zUTBGQlF5SXNJbVpwYkdVaU9pSjNkM2N2Y0dGamEzTXZkR1YwYUdWeUwzTnlZeTlxY3k5dFlYSnJRWFIwWVdOb2JXVnVkQzVxY3lJc0luTnZkWEpqWlhORGIyNTBaVzUwSWpwYklpOHFJR2RzYjJKaGJITWdWR1YwYUdWeUlDb3ZYRzVjYmxSbGRHaGxjaTV0YjJSMWJHVnpMbkIxYzJnb2UxeHVJQ0JwYm1sMGFXRnNhWHBsS0NrZ2UxeHVJQ0FnSUhSb2FYTXViV0Z5YTJWeWN5QTlJSHQ5TzF4dVhHNGdJQ0FnV3lkMFlYSm5aWFFuTENBblpXeGxiV1Z1ZENkZExtWnZja1ZoWTJnb2RIbHdaU0E5UGlCN1hHNGdJQ0FnSUNCamIyNXpkQ0JsYkNBOUlHUnZZM1Z0Wlc1MExtTnlaV0YwWlVWc1pXMWxiblFvSjJScGRpY3BPMXh1SUNBZ0lDQWdaV3d1WTJ4aGMzTk9ZVzFsSUQwZ2RHaHBjeTVuWlhSRGJHRnpjeWhnSkhzZ2RIbHdaU0I5TFcxaGNtdGxjbUFwTzF4dVhHNGdJQ0FnSUNCamIyNXpkQ0JrYjNRZ1BTQmtiMk4xYldWdWRDNWpjbVZoZEdWRmJHVnRaVzUwS0Nka2FYWW5LVHRjYmlBZ0lDQWdJR1J2ZEM1amJHRnpjMDVoYldVZ1BTQjBhR2x6TG1kbGRFTnNZWE56S0NkdFlYSnJaWEl0Wkc5MEp5azdYRzRnSUNBZ0lDQmxiQzVoY0hCbGJtUkRhR2xzWkNoa2IzUXBPMXh1WEc0Z0lDQWdJQ0IwYUdselczUjVjR1ZkTG1Gd2NHVnVaRU5vYVd4a0tHVnNLVHRjYmx4dUlDQWdJQ0FnZEdocGN5NXRZWEpyWlhKelczUjVjR1ZkSUQwZ2UyUnZkQ3dnWld4OU8xeHVJQ0FnSUgwcE8xeHVJQ0I5TEZ4dVhHNGdJSEJ2YzJsMGFXOXVLSHR0WVc1MVlXeFBabVp6WlhRc0lHMWhiblZoYkZSaGNtZGxkRTltWm5ObGRIMHBJSHRjYmlBZ0lDQmpiMjV6ZENCdlptWnpaWFJ6SUQwZ2UxeHVJQ0FnSUNBZ1pXeGxiV1Z1ZERvZ2JXRnVkV0ZzVDJabWMyVjBMRnh1SUNBZ0lDQWdkR0Z5WjJWME9pQnRZVzUxWVd4VVlYSm5aWFJQWm1aelpYUmNiaUFnSUNCOU8xeHVYRzRnSUNBZ1ptOXlJQ2hzWlhRZ2RIbHdaU0JwYmlCdlptWnpaWFJ6S1NCN1hHNGdJQ0FnSUNCamIyNXpkQ0J2Wm1aelpYUWdQU0J2Wm1aelpYUnpXM1I1Y0dWZE8xeHVJQ0FnSUNBZ1ptOXlJQ2hzWlhRZ2MybGtaU0JwYmlCdlptWnpaWFFwSUh0Y2JpQWdJQ0FnSUNBZ2JHVjBJSFpoYkNBOUlHOW1abk5sZEZ0emFXUmxYVHRjYmlBZ0lDQWdJQ0FnWTI5dWMzUWdibTkwVTNSeWFXNW5JRDBnZEhsd1pXOW1JSFpoYkNBaFBUMGdKM04wY21sdVp5YzdYRzRnSUNBZ0lDQWdJR2xtSUNodWIzUlRkSEpwYm1jZ2ZIeGNiaUFnSUNBZ0lDQWdJQ0FnSUhaaGJDNXBibVJsZUU5bUtDY2xKeWtnUFQwOUlDMHhJQ1ltWEc0Z0lDQWdJQ0FnSUNBZ0lDQjJZV3d1YVc1a1pYaFBaaWduY0hnbktTQTlQVDBnTFRFcElIdGNiaUFnSUNBZ0lDQWdJQ0IyWVd3Z0t6MGdKM0I0Snp0Y2JpQWdJQ0FnSUNBZ2ZWeHVYRzRnSUNBZ0lDQWdJR2xtSUNoMGFHbHpMbTFoY210bGNuTmJkSGx3WlYwdVpHOTBMbk4wZVd4bFczTnBaR1ZkSUNFOVBTQjJZV3dwSUh0Y2JpQWdJQ0FnSUNBZ0lDQjBhR2x6TG0xaGNtdGxjbk5iZEhsd1pWMHVaRzkwTG5OMGVXeGxXM05wWkdWZElEMGdkbUZzTzF4dUlDQWdJQ0FnSUNCOVhHNGdJQ0FnSUNCOVhHNGdJQ0FnZlZ4dVhHNGdJQ0FnY21WMGRYSnVJSFJ5ZFdVN1hHNGdJSDFjYm4wcE8xeHVJbDE5XG4iXX0=