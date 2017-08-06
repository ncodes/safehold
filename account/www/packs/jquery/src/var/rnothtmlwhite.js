define(function () {
    "use strict";
    // Only count HTML whitespace
    // Other whitespace should count in values
    // https://html.spec.whatwg.org/multipage/infrastructure.html#space-character
    return (/[^\x20\t\r\n\f]+/g);
});
//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbInd3dy9wYWNrcy9qcXVlcnkvc3JjL3Zhci9ybm90aHRtbHdoaXRlLmpzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUFBLE1BQU0sQ0FBRTtJQUNQLFlBQVksQ0FBQztJQUViLDZCQUE2QjtJQUM3QiwwQ0FBMEM7SUFDMUMsNkVBQTZFO0lBQzdFLE1BQU0sQ0FBQyxDQUFFLG1CQUFtQixDQUFFLENBQUM7QUFDaEMsQ0FBQyxDQUFFLENBQUMiLCJmaWxlIjoid3d3L3BhY2tzL2pxdWVyeS9zcmMvdmFyL3Jub3RodG1sd2hpdGUuanMiLCJzb3VyY2VzQ29udGVudCI6WyJkZWZpbmUoIGZ1bmN0aW9uKCkge1xuXHRcInVzZSBzdHJpY3RcIjtcblxuXHQvLyBPbmx5IGNvdW50IEhUTUwgd2hpdGVzcGFjZVxuXHQvLyBPdGhlciB3aGl0ZXNwYWNlIHNob3VsZCBjb3VudCBpbiB2YWx1ZXNcblx0Ly8gaHR0cHM6Ly9odG1sLnNwZWMud2hhdHdnLm9yZy9tdWx0aXBhZ2UvaW5mcmFzdHJ1Y3R1cmUuaHRtbCNzcGFjZS1jaGFyYWN0ZXJcblx0cmV0dXJuICggL1teXFx4MjBcXHRcXHJcXG5cXGZdKy9nICk7XG59ICk7XG4iXX0=

//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbInd3dy9wYWNrcy9qcXVlcnkvc3JjL3Zhci9ybm90aHRtbHdoaXRlLmpzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUFBLE1BQU0sQ0FBQztJQUNILFlBQVksQ0FBQztJQUNiLDZCQUE2QjtJQUM3QiwwQ0FBMEM7SUFDMUMsNkVBQTZFO0lBQzdFLE1BQU0sQ0FBQyxDQUFDLG1CQUFtQixDQUFDLENBQUM7QUFDakMsQ0FBQyxDQUFDLENBQUM7QUFFSCxteEJBQW14QiIsImZpbGUiOiJ3d3cvcGFja3MvanF1ZXJ5L3NyYy92YXIvcm5vdGh0bWx3aGl0ZS5qcyIsInNvdXJjZXNDb250ZW50IjpbImRlZmluZShmdW5jdGlvbiAoKSB7XG4gICAgXCJ1c2Ugc3RyaWN0XCI7XG4gICAgLy8gT25seSBjb3VudCBIVE1MIHdoaXRlc3BhY2VcbiAgICAvLyBPdGhlciB3aGl0ZXNwYWNlIHNob3VsZCBjb3VudCBpbiB2YWx1ZXNcbiAgICAvLyBodHRwczovL2h0bWwuc3BlYy53aGF0d2cub3JnL211bHRpcGFnZS9pbmZyYXN0cnVjdHVyZS5odG1sI3NwYWNlLWNoYXJhY3RlclxuICAgIHJldHVybiAoL1teXFx4MjBcXHRcXHJcXG5cXGZdKy9nKTtcbn0pO1xuXG4vLyMgc291cmNlTWFwcGluZ1VSTD1kYXRhOmFwcGxpY2F0aW9uL2pzb247Y2hhcnNldD11dGY4O2Jhc2U2NCxleUoyWlhKemFXOXVJam96TENKemIzVnlZMlZ6SWpwYkluZDNkeTl3WVdOcmN5OXFjWFZsY25rdmMzSmpMM1poY2k5eWJtOTBhSFJ0Ykhkb2FYUmxMbXB6SWwwc0ltNWhiV1Z6SWpwYlhTd2liV0Z3Y0dsdVozTWlPaUpCUVVGQkxFMUJRVTBzUTBGQlJUdEpRVU5RTEZsQlFWa3NRMEZCUXp0SlFVVmlMRFpDUVVFMlFqdEpRVU0zUWl3d1EwRkJNRU03U1VGRE1VTXNOa1ZCUVRaRk8wbEJRemRGTEUxQlFVMHNRMEZCUXl4RFFVRkZMRzFDUVVGdFFpeERRVUZGTEVOQlFVTTdRVUZEYUVNc1EwRkJReXhEUVVGRkxFTkJRVU1pTENKbWFXeGxJam9pZDNkM0wzQmhZMnR6TDJweGRXVnllUzl6Y21NdmRtRnlMM0p1YjNSb2RHMXNkMmhwZEdVdWFuTWlMQ0p6YjNWeVkyVnpRMjl1ZEdWdWRDSTZXeUprWldacGJtVW9JR1oxYm1OMGFXOXVLQ2tnZTF4dVhIUmNJblZ6WlNCemRISnBZM1JjSWp0Y2JseHVYSFF2THlCUGJteDVJR052ZFc1MElFaFVUVXdnZDJocGRHVnpjR0ZqWlZ4dVhIUXZMeUJQZEdobGNpQjNhR2wwWlhOd1lXTmxJSE5vYjNWc1pDQmpiM1Z1ZENCcGJpQjJZV3gxWlhOY2JseDBMeThnYUhSMGNITTZMeTlvZEcxc0xuTndaV011ZDJoaGRIZG5MbTl5Wnk5dGRXeDBhWEJoWjJVdmFXNW1jbUZ6ZEhKMVkzUjFjbVV1YUhSdGJDTnpjR0ZqWlMxamFHRnlZV04wWlhKY2JseDBjbVYwZFhKdUlDZ2dMMXRlWEZ4NE1qQmNYSFJjWEhKY1hHNWNYR1pkS3k5bklDazdYRzU5SUNrN1hHNGlYWDA9XG4iXX0=