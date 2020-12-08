// Environment view
export default {
    props: ['state', 'environmentUUID'],
    computed: {
        environment: function() {
            if(this.state.environments && this.state.environments.length)
            {
                return this.state.environments.find(e => e.UUID == this.environmentUUID);
            }
            return {};
        },
    },
    template: `
        <div class="environment-container">

            <h1><i class="fas fa-cloud-sun"></i> {{ environment.Name }}</h1>

            <div class="row">
                <!-- Environment Meta -->
                <div class="col-sm">
                    <ganjapp-environment-meta :environment="environment"></ganjapp-environment-meta>
                </div>
                
                <!-- Environment Status -->
                <div class="col-sm">
                    <ganjapp-environment-status :environment="environment"></ganjapp-environment-status>
                </div>
            </div>

            <!-- Environment Events -->
            <ganjapp-events v-if="environment.Events.length" :events="environment.Events" :toasts="state.toasts" :limit="10"></ganjapp-events>

            <!-- Environment Extended Data -->
            <ganjapp-environment-extended-data :environment="environment"></ganjapp-environment-extended-data>

        </div>
    `
}