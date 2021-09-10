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

    <!-- This lengthy to replace javascript functionality -->
    {{ if eq $.help "windows" }}
      <div class="columns is-centered">
        <div class="column">
          <div class="tabs is-boxed is-medium is-centered">
            <ul>
              <li><a href="?help=mac">Mac</a></li>
              <li class="is-active"><a href="?help=windows">Windows</a></li>
              <li><a href="?help=linux">Linux</a></li>
              <li><a href="?help=android">Android</a></li>
              <li><a href="?help=ios">iOS</a></li>
            </ul>
          </div>
          <div>
            Help for Windows
          </div>
        </div>
      </div> <!--columns-->
    {{ else if eq $.help "linux" }}
      <div class="columns is-centered">
        <div class="column">
          <div class="tabs is-boxed is-medium is-centered">
            <ul>
              <li><a href="?help=mac">Mac</a></li>
              <li><a href="?help=windows">Windows</a></li>
              <li class="is-active"><a href="?help=linux">Linux</a></li>
              <li><a href="?help=android">Android</a></li>
              <li><a href="?help=ios">iOS</a></li>
            </ul>
          </div>
          <div>
            Help for Linux
          </div>
        </div>
      </div> <!--columns-->
    {{ else if eq $.help "android" }}
      <div class="columns is-centered">
        <div class="column">
          <div class="tabs is-boxed is-medium is-centered">
            <ul>
              <li><a href="?help=mac">Mac</a></li>
              <li><a href="?help=windows">Windows</a></li>
              <li><a href="?help=linux">Linux</a></li>
              <li class="is-active"><a href="?help=android">Android</a></li>
              <li><a href="?help=ios">iOS</a></li>
            </ul>
          </div>
          <div>
            Help for Android
          </div>
        </div>
      </div> <!--columns-->
    {{ else if eq $.help "ios" }}
      <div class="columns is-centered">
        <div class="column">
          <div class="tabs is-boxed is-medium is-centered">
            <ul>
              <li><a href="?help=mac">Mac</a></li>
              <li><a href="?help=windows">Windows</a></li>
              <li><a href="?help=linux">Linux</a></li>
              <li><a href="?help=android">Android</a></li>
              <li class="is-active"><a href="?help=ios">iOS</a></li>
            </ul>
          </div>
          <div>
            Help for iOS
          </div>
        </div>
      </div> <!--columns-->
    {{ else }}
      <div class="columns is-centered">
        <div class="column">
          <div class="tabs is-boxed is-medium is-centered">
            <ul>
              <li class="is-active"><a href="?help=mac">Mac</a></li>
              <li><a href="?help=windows">Windows</a></li>
              <li><a href="?help=linux">Linux</a></li>
              <li><a href="?help=android">Android</a></li>
              <li><a href="?help=ios">iOS</a></li>
            </ul>
          </div>
          <div>
            Help for Mac
          </div>
        </div>
      </div> <!--columns-->
    {{ end }}
      </section>
    </div> <!--container-->
  </body>
</html>
