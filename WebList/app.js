var app = new Vue({
    el: '#app',
    data: {
      message: 'Привет, Vue!',
      InpFormMessage: 'Введите сообщение',
      btnName: 'Сгенерировать',
      OutpKeyMess: 'Ключ',
      OutpURLMess: 'Адрес',
      Message: ""
    },
    methods:{
      CheckText: function(){
        console.log('Chekau', this.Message)
              // GET /someUrl
        this.$http.post('http://localhost:9090//makeText', {MSG : this.Message}).then(res => {
            console.log('Generation message', res)

        }, res => {
            console.log("SendError", res)
        });
      }
    }
  })