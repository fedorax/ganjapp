// Toast container component
export default {
    name: 'GanjappToastContainerComponent',
    props: ['toasts'],
    template: `
        <teleport to=".ganjapp-toast-container">
            <div class="ganjapp-toast-wrapper fixed-top">
                <ganjapp-toast v-for="t in toasts" v-bind:toast="t" v-bind:toasts="toasts"></ganjapp-toast>
            </div>
        </teleport>
    `,
    data() {
        return {};
    }
}