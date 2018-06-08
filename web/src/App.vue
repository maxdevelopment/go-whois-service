<template>
  <div id="app">
    <b-table striped hover :items="usersTable" :fields="fields"></b-table>
  </div>
</template>

<script>
  import Routes from '~/static/routes'

  export default {
    name: 'app',
    data: () => ({
      id: null,
      ws: null,
      usersTable: [],
      fields: [ 'id', 'remote_addr', 'city', 'country', 'cached', 'link' ],
    }),
    beforeMount() {
      this.id = Math.random().toString(36).substr(2, 9)
      this.ws = new WebSocket(Routes.socketPath(this.id))
      this.ws.addEventListener('message', (e) => {
        let result = []
        let data = JSON.parse(e.data)
        for (let val of Object.values(data)) {
          result.push(val)
        }
        this.usersTable = result
      })
    }
  }
</script>

<style>
  #app {
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
    color: #2c3e50;
    margin-top: 60px;
  }

  h1, h2 {
    font-weight: normal;
  }

  ul {
    list-style-type: none;
    padding: 0;
  }

  li {
    display: inline-block;
    margin: 0 10px;
  }

  a {
    color: #42b983;
  }
</style>
