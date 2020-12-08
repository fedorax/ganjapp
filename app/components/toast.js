// Toast component
export default {
    name: 'GanjappToastComponent',
    props: ['toasts', 'toast'],
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

        removeToast: function() {
            for(var i = 0; i < this.toasts.length; i++)
            {
                var t = this.toasts[i];
                if(t.id == this.toast.id)
                {
                    this.toasts.splice(i, 1);
                    break;
                }
            }
        }

    },
    template: `
        <div class="toast" data-autohide="false" role="alert" aria-live="assertive" aria-atomic="true">
            <div class="toast-header">
                <img src="/icons/favicon-16x16.png" class="mr-2">
                <strong class="mr-auto"><span class="d-none d-sm-inline">{{ titleCase(toast.object.type) }}@</span>{{ kebabCaseToTitleCase(toast.event) }}</strong>
                <span class="badge ml-2" v-bind:class="alertClass(toast.severity)" role="alert">{{ titleCase(toast.severity) }}</span>
                <button type="button" class="ml-2 mb-1 close" data-dismiss="toast" v-on:click="removeToast()" aria-label="Close">
                    <span aria-hidden="true">&times;</span>
                </button>
            </div>
            <div class="toast-body">
                {{ toast.message }}
            </div>
        </div>
    `,
    mounted() {
        $('.toast').toast('show');
    },
    beforeUpdate() {
        $('.toast').toast('show');
    },
    update() {
        $('.toast').toast('show');
    },
    data() {
        return {};
    }
}