<html>
  {{ template "header" . }}
  <body>
    <div class="container is-fluid">

      {{ template "display_errors" . }}

      <section id="devices" class="section">
      <div class="columns is-centered">
        <div class="column has-background-light">

          <div>
            <a class="button is-primary" href="/devices">
              <span class="icon">
                <i class="fas fa-plus"></i>
              </span>
              <span>Add device</span>
            </a>
          </div>
          <br>
          <table class="table is-striped is-fullwidth is-hoverable">
            <thead>
              <tr>
                <th>Name</th>
                <th>Public key</th>
                <th>Status</th>
                <th>Statistics</th>
                <th>Actions</th>
              </tr>
            </thead>

            <tbody>
              {{ range $i, $d := $.devices }}
              <tr>
                <td>{{ $d.dev.Name }}</td>
                <td>{{ $d.dev.PrivateKey.PublicKey }}</td>
                <td>{{ $d.stat }}</td>
                <td>
                  <div>
                  Last seen: {{ $d.peer.LastHandshakeTime }}
                  </div>
                  <div>
                  Received bytes: {{ $d.peer.ReceiveBytes }}
                  </div>
                  <div>
                  Transmitted bytes: {{ $d.peer.TransmitBytes }}
                  </div>
                </td>
                <td>
                  <div>
                    <a class="button is-primary" href="/devices/{{ $d.dev.Id }}/install">
                      <span class="icon">
                        <i class="fas fa-tools"></i>
                      </span>
                      <span>
                        Install
                      </span>
                    </a>
                    <a class="button is-danger" href="/devices/{{ $d.dev.Id }}/delete" onclick="return confirmDelete()">
                      <span class="icon">
                        <i class="fas fa-trash-alt"></i>
                      </span>
                      <span>
                        Delete
                      </span>
                    </a>
                  </div>
                </td>
              </tr>
              {{ end }}
            </tbody>
          </table>

        </div>
      </div> <!--columns-->
      </section>
    </div> <!--container-->

<script>
function confirmDelete() {
  var ans = confirm("Do you really want to delete this device?");
  return ans;
}
</script>

  </body>
</html>
