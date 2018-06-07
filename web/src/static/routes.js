const wsSchema = window.location.protocol === 'http:' ? 'ws:' : 'wss:'
const baseSocket = window.location.hostname === '127.0.0.1' ? `${wsSchema}//127.0.0.1:8000` : `${wsSchema}//${window.location.hostname}`

export default {
  socketPath: (id) => (
    `${baseSocket}/join/${id}`
  )
}
