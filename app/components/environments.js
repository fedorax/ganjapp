// Environments component
export default {
    name: 'GanjappEnvironmentsComponent',
    props: ['environments'],
    template: `
        <h3 v-if="environments.length"><i class="fas fa-cloud-sun"></i> Environments</h3>
        <table class="table table-hover table-sm" v-if="environments.length">
            <thead class="bg-success">
                <tr>
                    <th scope="col" title="Environment"><i class="fas fa-cloud-sun"></i><span class="d-none d-sm-inline"><br />Environment</span></th>
                    <th scope="col" title="Temperature"><i class="fas fa-thermometer-half"></i><span class="d-none d-sm-inline"><br />Temperature</span></th>
                    <th scope="col" title="Humidity"><i class="fas fa-tint"></i><span class="d-none d-sm-inline"><br />Humidity</span></th>
                    <th scope="col" title="Lighting"><i class="fas fa-lightbulb"></i><span class="d-none d-sm-inline"><br />Lighting</span></th>
                    <th scope="col" class="d-none d-sm-table-cell" title="Trees"><i class="fas fa-seedling"></i><span class="d-none d-sm-inline"><br />Trees</span></th>
                    <th scope="col" class="d-none d-sm-table-cell" title="Shrooms"><i class="fas fa-long-arrow-alt-up"></i><span class="d-none d-sm-inline"><br />Shrooms</span></th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="e in environments" style="transform: rotate(0);">
                    <th scope="row">
                        <router-link class="text-decoration-none text-success stretched-link" :to="{ path: '/environment/' + e.UUID + '/'}">{{ e.Name }}</router-link>
                    </th>
                    <td>
                        <span v-if="e.Status.Temperature != null">{{ e.Status.Temperature }}&deg;C</span>
                        <span v-if="e.Status.Temperature == null">-</span>
                    </td>
                    <td>
                        <span v-if="e.Status.Humidity != null">{{ e.Status.Humidity }}%</span>
                        <span v-if="e.Status.Humidity == null">-</span>
                    </td>
                    <td>
                        <span v-if="e.Status.LightsOn != null">{{ e.Status.LightsOn ? "On" : "Off" }}</span>
                        <span v-if="e.Status.LightsOn == null">-</span>
                    </td>
                    <td class="d-none d-sm-table-cell">{{ e.Trees.length }}</td>
                    <td class="d-none d-sm-table-cell">{{ e.Shrooms.length }}</td>
                </tr>
            </tbody>
        </table>
    `,
    data() {
        return {};
    }
}