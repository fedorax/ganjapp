// Events component
export default {
    name: 'GanjappEventsComponent',
    props: ['events', 'toasts', 'limit'],
    emits: ['create-toast'],
    methods: {

        lowerCase: function(s) {
            return s.toLowerCase();
        },

        titleCase: function(s) {
            return s.toLowerCase().split(' ').map(function(w) {
                return w.replace(w[0], w[0].toUpperCase());
            }).join(' ');
        },

        kebabCaseToTitleCase: function(s) {
            return this.titleCase(s.replaceAll("-", " "));
        },

        alertClass: function(s) {
            s = this.lowerCase(s);
            return {
                'badge-info': s == 'info',
                'badge-warning': s == 'warning',
                'badge-danger': s == 'error',
            };
        },

        parseTime: function(t) {
            var d = new Date(t);
            return d.toLocaleString();
        },

        createToast: function(e) {
            this.removeToast(e.ID);
            this.toasts.push({
              id: e.ID,
              timestamp: e.CreatedAt,
              severity: e.Severity,
              event: e.Event,
              object: {
                id: e.ObjectID,
                type: e.ObjectType,
                user: e.UserID,
              },
              message: e.Message,
              data: e.Data,
            });
          },
    
          removeToast: function(e) {
            for(var i = 0; i < this.toasts.length; i++)
            {
              var t = this.toasts[i];
              if(t.id == e)
              {
                this.toasts.splice(i, 1);
                break;
              }
            }
          },

    },
    template: `
        <h3 v-if="events.length"><i class="fas fa-bell"></i> Events</h3>
        <table class="table table-hover table-sm" v-if="events.length">
            <thead class="bg-success">
                <tr>
                    <th scope="col" class="d-none d-sm-table-cell" title="Timestamp"><i class="fas fa-clock"></i><br />Timestamp</th>
                    <th scope="col" title="Severity"><i class="fas fa-exclamation"></i><span class="d-none d-sm-inline"><br />Severity</span></th>
                    <th scope="col" class="d-none d-sm-table-cell" title="Component"><i class="fas fa-server"></i><span class="d-none d-sm-inline"><br />Component</span></th>
                    <th scope="col" class="d-none d-sm-table-cell" title="Event"><i class="fas fa-bell"></i><span class="d-none d-sm-inline"><br />Event</span></th>
                    <th scope="col" title="Message"><i class="fas fa-database"></i><span class="d-none d-sm-inline"><br />Message</span></th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="e in (limit ? events.slice(0, (Number.isInteger(limit) ? limit : 10)) : events)" v-on:click="createToast(e)">
                    <td class="align-middle d-none d-sm-table-cell">{{ parseTime(e.CreatedAt) }}</td>
                    <th scope="row">
                        <div class="badge" v-bind:class="alertClass(e.Severity)" role="alert">{{ titleCase(e.Severity) }}</div>
                    </th>
                    <td class="align-middle d-none d-sm-table-cell">{{ titleCase(e.ObjectType) }}</td>
                    <td class="align-middle d-none d-sm-table-cell">{{ kebabCaseToTitleCase(e.Event) }}</td>
                    <td class="align-middle">{{ e.Message }}</td>
                </tr>
            </tbody>
        </table>
    `,
    data() {
        return {};
    }
}