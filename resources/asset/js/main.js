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
                let port = location.port ? `:${location.port}` : ''
                this.shortUrl = location.protocol + '//' + location.hostname + port + '/link/' + res.data.data
                this.value = this.orgUrl
                this.makeQrcode()
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
        setLogo() {
            (function($){
                function injector(t, splitter, klass, after) {
                    var a = t.text().split(splitter), inject = '';
                    if (a.length) {
                        $(a).each(function(i, item) {
                            inject += '<span class="'+klass+(i+1)+'">'+item+'</span>'+after;
                        });
                        t.empty().append(inject);
                    }
                }

                var methods = {
                    init : function() {

                        return this.each(function() {
                            injector($(this), '', 'char', '');
                        });

                    },

                    words : function() {

                        return this.each(function() {
                            injector($(this), ' ', 'word', ' ');
                        });

                    },

                    lines : function() {

                        return this.each(function() {
                            var r = "eefec303079ad17405c889e092e105b0";
                            // Because it's hard to split a <br/> tag consistently across browsers,
                            // (*ahem* IE *ahem*), we replaces all <br/> instances with an md5 hash
                            // (of the word "split").  If you're trying to use this plugin on that
                            // md5 hash string, it will fail because you're being ridiculous.
                            injector($(this).children("br").replaceWith(r).end(), r, 'line', '');
                        });

                    }
                };

                $.fn.lettering = function( method ) {
                    // Method calling logic
                    if ( method && methods[method] ) {
                        return methods[ method ].apply( this, [].slice.call( arguments, 1 ));
                    } else if ( method === 'letters' || ! method ) {
                        return methods.init.apply( this, [].slice.call( arguments, 0 ) ); // always pass an array
                    }
                    $.error( 'Method ' +  method + ' does not exist on jQuery.lettering' );
                    return this;
                };

            })(jQuery);
            $(document).ready(function() {
                $("h1.logo").lettering();
            });
        }
    },
    mounted() {
        this.setToastr()
    }
}

Vue.createApp(app).mount('#app')
