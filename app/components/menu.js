// Menu component
const axios = require('axios').default;
export default {
    name: 'GanjappMenuComponent',
    props: [
        'state',
    ],
    methods: {
        isCurrentRoute: function(n) {
            return this.$route.name == n;
        },
        createEnvironment: function() {
            axios.post('/api/create/environment', {
                name: this.newEnvironmentName,
            })
            .catch(function (e) {
                console.warn(error);
            });
        },
    },
    template: `
        <teleport to="#navbarContainer">
            <ul class="navbar-nav">
                <li class="nav-item">
                    <router-link to="/" class="nav-link" :class="{ 'active': isCurrentRoute('dashboard') }">Dashboard</router-link>
                </li>
                <li class="nav-item dropdown" :class="{ 'active': isCurrentRoute('environment') }">
                    <a class="nav-link dropdown-toggle" href="#" id="environmentsDropdownMenuLink" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">Environments</a>
                    <div class="dropdown-menu" aria-labelledby="environmentsDropdownMenuLink">
                        <h6 class="dropdown-header" v-if="state.environments.length">Environments</h6>
                        <router-link v-if="state.environments.length" v-for="e in state.environments" :to="{ path: '/environment/' + e.UUID + '/'}" class="dropdown-item" onclick="$('#environmentsDropdownMenuLink').dropdown('hide')">{{ e.Name }}</router-link>
                        <div class="dropdown-divider" v-if="state.environments.length"></div>
                        <h6 class="dropdown-header">Create Environment</h6>
                        <form name="createEnvironment" autocomplete="off" class="px-4 py-3" method="post" action="/api/create/environment/">
                            <input type="hidden" autocomplete="false">
                            <div class="form-group">
                                <input type="text" class="form-control form-control-sm" id="newEnvironmentName" v-model="newEnvironmentName" placeholder="Name" required>
                            </div>
                            <button type="button" class="btn btn-sm btn-success btn-block" :disabled="!newEnvironmentName" @click="createEnvironment()">Create</button>
                        </form>
                    </div>
                </li>
                <li class="nav-item">
                    <router-link to="/events" class="nav-link" :class="{ 'active': isCurrentRoute('events') }">Events</router-link>
                </li>
            </ul>
        </teleport>
    `,
    data() {
        return {
            newEnvironmentName: null,
        }
    },
}