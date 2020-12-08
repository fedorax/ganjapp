// Environment Extended Data component
export default {
    name: 'GanjappEnvironmentExtendedDataComponent',
    props: ['environment'],
    methods: {
        addNewData: function() {
            this.environment.ExtendedData.push({
                Key: "set-me-" + this.environment.ExtendedData.length,
                Value: "0",
            });
        },

        removeData: function(index) {
            this.environment.ExtendedData.splice(index, 1);
        }

    },
    template: `
        <h3><i class="fas fa-plus"></i> Custom Data <button class="btn btn-sm btn-success" @click="addNewData()">New</button></h3>
        <table class="table" v-if="environment.ExtendedData.length">
            <tbody>
                <tr v-for="(d, i) in environment.ExtendedData" :title="d.Key">
                    <td class="text-right">
                        <div class="input-group">
                            <div class="input-group-prepend">
                                <span class="input-group-text" :id="'ecp-icon' + i"><i class="fas fa-key"></i></span>
                            </div>
                            <input class="form-control" type="text" v-model="d.Key" placeholder="Custom Data Key" aria-label="Custom Data Key" :aria-describedby="'ecp-icon' + i" required>
                        </div>
                    </td>
                    <td class="text-left">
                        <div class="input-group">
                            <div class="input-group-prepend">
                                <span class="input-group-text" :id="'ecv-icon' + i"><i class="fas fa-database"></i></span>
                            </div>
                            <input class="form-control" type="text" v-model="d.Value" placeholder="Custom Data Value" aria-label="Custom Data Value" :aria-describedby="'ecv-icon' + i" required>
                        </div>
                    </td>
                    <td class="text-center">
                        <div class="btn-group" role="group" aria-label="Custom Data Functions">
                            <button class="btn btn-success" title="Save"><i class="fas fa-save"></i></button>
                            <button class="btn btn-danger" title="Delete" @click="removeData(i)"><i class="fas fa-times"></i></button>
                        </div>
                    </td>
                </tr>
            </tbody>
        </table>
    `
}