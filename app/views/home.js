// Home view
export default {
    props: ['state'],
    template: `
        <ganjapp-environments :environments="state.environments"></ganjapp-environments>

        <ganjapp-events :events="state.events" :toasts="state.toasts" :limit="10"></ganjapp-events>
    `
}