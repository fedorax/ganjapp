// Environment Meta component
export default {
    name: 'GanjappEnvironmentMetaComponent',
    props: ['environment'],
    methods: {
        parseTime: function(t) {
            var d = new Date(t);
            return d.toLocaleString();
        },
    },
    template: `
        <h3><i class="fas fa-cogs"></i> Meta</h3>
        <table class="table table-sm table-hover">
            <tbody>
                <tr>
                    <td><i class="fas fa-id-card-alt"></i><span class="d-none d-sm-inline">&nbsp;<strong>ID</strong></td>
                    <td>{{ environment.UUID }}</td>
                </tr>
                <tr>
                    <td><i class="fas fa-plus"></i><span class="d-none d-sm-inline">&nbsp;<strong>Created</strong></td>
                    <td>{{ parseTime(environment.CreatedAt) }}</td>
                </tr>
                <tr>
                    <td><i class="fas fa-pen"></i><span class="d-none d-sm-inline">&nbsp;<strong>Updated</strong></td>
                    <td>{{ parseTime(environment.UpdatedAt) }}</td>
                </tr>
                <tr>
                    <td>
                        <i class="fas fa-comment"></i><span class="d-none d-sm-inline">&nbsp;<strong>Comments</strong></span>
                        <br />
                        <button class="btn btn-sm btn-success"><i class="fas fa-save"></i></button>
                    </td>
                    <td>
                        <textarea class="form-control" placeholder="Environment Comments" v-model="environment.Comments"></textarea>
                    </td>
                </tr>
            </tbody>
        </table>
    `
}