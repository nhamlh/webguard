<html>
  {{ template "header" }}
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
                    <a class="button is-primary" href="/devices/{{ .id }}/download">
                      <span class="icon">
                        <i class="fas fa-download"></i>
                      </span>
                    </a>
                    <a class="button is-danger" href="/devices/{{ .id }}/delete">
                      <span class="icon">
                        <i class="fas fa-trash-alt"></i>
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

      <div class="is-divider"></div>

      <section id="help" class="section">
      <div class="columns is-centered">
        <div class="column has-text-centered">
          <div class="title">
            Help
          </div>
        </div>
      </div> <!--columns-->

      <div class="columns is-centered">
        <div class="column">
          <div class="tabs is-boxed is-medium is-centered">
            <ul>
              <li class="is-active"><a>Mac</a></li>
              <li><a>Windows</a></li>
              <li><a>Linux</a></li>
              <li><a>Android</a></li>
              <li><a>iOS</a></li>
            </ul>
          </div>
          <div>
Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.
          </div>
        </div>
      </div> <!--columns-->
      </section>
    </div> <!--container-->
  </body>
</html>
