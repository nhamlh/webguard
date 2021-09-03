<html>
  {{ template "header" }}
  <body>
    <div class="container">

      {{ template "display_errors" . }}

      <br>

      <section class="section">
        <div><a href="/new_device" class="button">Add New Device</a></div>
        <br>
        <table class="table">
          <thead>
            <tr>
              <th>Id</th>
              <th>Name</th>
              <th>Public key</th>
              <th>Last seen</th>
              <th>Actions</th>
            </tr>
          </thead>

          <tbody>
            {{ range $.devices }}
            <tr>
              <td>{{ .id }}</td>
              <td>{{ .name }}</td>
              <td>{{ .pubkey }}</td>
              <td>{{ .lastSeen }}</td>
              <td><a class="button is-primary" href="/devices/{{ .id }}/download">Download</a> <a class="button is-danger" href="/devices/{{ .id }}/delete">Delete</a></td>
            </tr>
            {{ end }}
          </tbody>
        </table>
      </section>

      <section class="section">
        <p class="header">Help</p>
        <div class="tabs is-boxed is-medium is-centered">
          <ul>
            <li class="is-active"><a>Mac</a></li>
            <li><a>Windows</a></li>
            <li><a>Linux</a></li>
            <li><a>Android</a></li>
            <li><a>iOS</a></li>
          </ul>
        </div>
      </section>

    </div> <!--container-->
  </body>
</html>
