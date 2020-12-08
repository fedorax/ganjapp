// Events view
export default {
    props: ['state'],
    template: `
        <ganjapp-events :events="state.events" :toasts="state.toasts"></ganjapp-events>
    `
}