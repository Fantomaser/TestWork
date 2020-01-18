var app = new Vue({
    el: '#app',
    data: {
      message: null,
      btnName: 'Сгенерировать',
      OutpKeyMess: 'Ключ',
      OutpURLMess: 'Адрес' 
    },
    methods:{
      CheckText: function(){
        console.log('Chekau', this.message)
              // GET /someUrl
        this.$http.post('http://localhost:8080/makeText', {msg : this.message}).then(res => {
            console.log('Generation message', res)

        }, res => {
            console.log("SendError!!!!!!", res)
        });
      }
    }
  })