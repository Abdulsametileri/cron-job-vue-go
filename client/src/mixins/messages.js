export default {
    methods: {
        async showErrorMessage(msg = "", title = "Error") {
            await this.$bvModal.msgBoxOk(msg, {
                title: title,
                headerBgVariant: 'danger',
                headerTextVariant: 'light',
                size: 'md',
                buttonSize: 'md',
                okVariant: 'danger',
                headerClass: 'p-2 border-bottom-0',
                bodyClass: 'modalCustomBody',
                footerClass: 'p-2 border-top-0',
            });
        },
        async showSucessMessage(msg = this.$t('operationSuccess'), title = 'Success') {
            await this.$bvModal.msgBoxOk(msg, {
                title: title,
                headerBgVariant: 'success',
                headerTextVariant: 'light',
                size: 'md',
                buttonSize: 'md',
                okVariant: 'success',
                headerClass: 'p-2 border-bottom-0',
                bodyClass: 'modalCustomBody',
                footerClass: 'p-2 border-top-0',
            });
        },
    }
}