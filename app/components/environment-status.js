// Environment Status component
export default {
    name: 'GanjappEnvironmentStatusComponent',
    props: ['environment'],
    template: `
        <h3><i class="fas fa-heartbeat"></i> Status</h3>
        <table class="table table-sm table-hover">
            <tbody>
                <tr title="Temperature">
                    <td><i class="fas fa-thermometer-half"></i><span class="d-none d-sm-inline">&nbsp;<strong>Temperature</strong></span></td>
                    <td>
                        <span v-if="environment.Status.Temperature != null">{{ environment.Status.Temperature }}&deg;C</span>
                        <span v-if="environment.Status.Temperature == null">-</span>
                    </td>
                </tr>
                <tr title="Humidity">
                    <td><i class="fas fa-tint"></i><span class="d-none d-sm-inline">&nbsp;<strong>Humidity</strong></span></td>
                    <td>
                        <span v-if="environment.Status.Humidity != null">{{ environment.Status.Humidity }}&deg;C</span>
                        <span v-if="environment.Status.Humidity == null">-</span>
                    </td>
                </tr>
                <tr title="Lighting">
                    <td><i class="fas fa-lightbulb"></i><span class="d-none d-sm-inline">&nbsp;<strong>Lighting</strong></span></td>
                    <td>
                        <span v-if="environment.Status.LightsOn != null">{{ environment.Status.LightsOn ? "On" : "Off" }}</span>
                        <span v-if="environment.Status.LightsOn == null">-</span>
                    </td>
                </tr>
            </tbody>
        </table>
    `
}