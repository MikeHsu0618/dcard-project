const app = {
    data() {
        return {
            orgUrl: "",
            shortUrl: "",
        }
    },
    methods: {
        async generateUrl() {
            // 驗證
            if (!this.orgUrl) {
                this.shortUrl = ""
                return
            }

            if (!this.orgUrl.includes(`https://`) && !this.orgUrl.includes(`http://`)) {
                this.shortUrl = ""
                toastr['error']('網址好像不太完整唷！', '錯誤')
                return
            }

            // 取得網址
            try {
                const res = await axios.post(`/link`, {org_url: this.orgUrl})
                port = location.port ? `:${location.port}` : ''
                this.shortUrl = location.hostname + port + '/link/' + res.data.data
            } catch (e) {
                this.shortUrl = ""
                toastr['error']('好像是個無效的網址唷！', '錯誤')
                console.log(e)
            }
        },
        copy() {
            navigator.clipboard.writeText(this.shortUrl)
                .then(() => {
                    toastr['success']('已複製到剪貼簿', '成功')
                }).catch(err => {
                toastr['error']('好像哪裡出錯囉！ 請稍候重試～', '錯誤')
            })
        },
        setToastr() {
            toastr.options = {
                "closeButton": true,
                "debug": false,
                "newestOnTop": false,
                "progressBar": true,
                "positionClass": "toast-top-right",
                "preventDuplicates": false,
                "onclick": null,
                "showDuration": "300",
                "hideDuration": "1000",
                "timeOut": "3000",
                "extendedTimeOut": "1000",
                "showEasing": "swing",
                "hideEasing": "linear",
                "showMethod": "fadeIn",
                "hideMethod": "fadeOut"
            }
        }
    },
    mounted() {
        this.setToastr()
    }
}


Vue.createApp(app).mount('#app')
