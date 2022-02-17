const app = {
    data() {
        return {
            orgUrl: "",
            shortUrl: "",
            value: "",
        }
    },
    methods: {
        async generateUrl() {
            // 驗證
            if (this.value === this.orgUrl) {
                return
            }

            if (!this.orgUrl) {
                this.shortUrl = ""
                this.value = ""
                return
            }

            if (!this.orgUrl.includes(`https://`) && !this.orgUrl.includes(`http://`)) {
                this.shortUrl = ""
                this.value = ""
                toastr['error']('網址好像不太完整唷！', '錯誤')
                return
            }

            // 取得網址
            try {
                const res = await axios.post(`/link`, {org_url: this.orgUrl})
                let port = location.port ? `:${location.port}` : ''
                this.shortUrl = location.protocol + '//' + location.hostname + port + '/link/' + res.data.data
                this.value = this.orgUrl
                this.makeQrcode()
            } catch (e) {
                this.shortUrl = ""
                this.value = ""
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
        },
        makeQrcode() {
            document.getElementById('qrcode').innerHTML = "";
            const qrcode = new QRCode(document.getElementById('qrcode'), {
                text: this.shortUrl,
                width: 128,
                height: 128,
                colorDark : '#000',
                colorLight : '#fff',
                correctLevel : QRCode.CorrectLevel.H
            });
            qrcode.makeCode(this.shortUrl);
        },
    },
    mounted() {
        this.setToastr()
    }
}

Vue.createApp(app).mount('#app')
