<html>
  {{ template "header" . }}
  <body>
    <div class="container is-fluid">

      {{ template "display_errors" . }}

      <section id="devices" class="section">
      <div class="columns is-centered">
        <div class="column has-background-light">

          <div>
            <a class="button is-primary" href="/new_device">
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
                <th>Last seen</th>
                <th>Actions</th>
              </tr>
            </thead>

            <tbody>
              {{ range $i, $d := $.devices }}
              <tr>
                <td>{{ $d.name }}</td>
                <td>{{ $d.pubkey }}</td>
                <td>{{ $d.lastSeen }}</td>
                <td>
                  <div>
                    <a class="button is-primary" href="/devices/{{ .id }}/install">
                      <span class="icon">
                        <i class="fas fa-tools"></i>
                      </span>
                      <span>
                        Install
                      </span>
                    </a>
                    <a class="button is-danger" href="/devices/{{ .id }}/delete" onclick="return confirmDelete()">
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
