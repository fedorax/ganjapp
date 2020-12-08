const Vue = require('./node_modules/vue/dist/vue.esm-bundler.js');
import { createWebHistory, createRouter } from "vue-router";
const axios = require('axios').default;

// Components
import GanjappMenuComponent from "./components/menu.js";
import GanjappToastContainerComponent from "./components/toast-container.js";
import GanjappToastComponent from "./components/toast.js";
import GanjappEventsComponent from "./components/events.js";
import GanjappEnvironmentsComponent from "./components/environments.js";
import GanjappEnvironmentMetaComponent from "./components/environment-meta.js";
import GanjappEnvironmentStatusComponent from "./components/environment-status.js";
import GanjappEnvironmentExtendedDataComponent from "./components/environment-extended-data.js";

// Views
const GanjappHomeView = require('./views/home.js').default;
const GanjappEventsView = require('./views/events.js').default;
const GanjappEnvironmentView = require('./views/environment.js').default;

const routes = [
  {
    path: "/",
    name: "dashboard",
    component: GanjappHomeView,
  },
  {
    path: "/environment/:environmentUUID",
    name: "environment",
    props: true,
    component: GanjappEnvironmentView,
  },
  {
    path: "/events",
    name: "events",
    component: GanjappEventsView,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

const ganjapp = Vue.createApp({
    el: '#app',
    data () {
      return {
        state: {
          events: [],
          environments: [],
          toasts: [],
        }
      }
    },
    methods: {

      createToast: function(e) {
        this.removeToast(e.ID);
        this.state.toasts.push({
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
        for(var i = 0; i < this.state.toasts.length; i++)
        {
          var t = this.state.toasts[i];
          if(t.id == e)
          {
            this.state.toasts.splice(i, 1);
            break;
          }
        }
      },

    },
    mounted () {
      
      // Fetch initial data
      axios.get('/api/events').then(response => (this.state.events = response.data));
      axios.get('/api/environments').then(response => (this.state.environments = response.data));

      // Connect to SSE for push updates
      var url = "/live/sse";
      var stream = new EventSource(url);

      stream.addEventListener("end", e => stream.close());

      stream.addEventListener("message", function(e)
      {
        console.log(e.data);
      });

      stream.addEventListener("event", e => {
        var event = JSON.parse(e.data);
        if(event) {
          this.state.events.splice(0, 0, event);
          this.createToast(event);
        }
      })

      stream.addEventListener("environment-update", e => {
        // Update the environment...
        var environment = JSON.parse(e.data);
        if(environment) {
          for(var i = 0; i < this.state.environments.length; i++) {
            if(this.state.environments[i].ID == environment.ID)
            {
              // Overwrite the existing environment...
              this.state.environments[i] = environment;
              return;
            }
          }

          // If we've made it this far, then there is no matching environment to update,
          // so just add it to the array...
          this.state.environments.push(environment);
          return;
        } else {
          console.error("Failed to parse environment-update data from server");
        }

      });

    }
});

ganjapp.use(router);

ganjapp.component("GanjappMenu", GanjappMenuComponent);
ganjapp.component("GanjappToastContainer", GanjappToastContainerComponent);
ganjapp.component("GanjappToast", GanjappToastComponent);
ganjapp.component("GanjappEvents", GanjappEventsComponent);
ganjapp.component("GanjappEnvironments", GanjappEnvironmentsComponent);
ganjapp.component("GanjappEnvironmentMeta", GanjappEnvironmentMetaComponent);
ganjapp.component("GanjappEnvironmentStatus", GanjappEnvironmentStatusComponent);
ganjapp.component("GanjappEnvironmentExtendedData", GanjappEnvironmentExtendedDataComponent);

const vm = ganjapp.mount("#app")