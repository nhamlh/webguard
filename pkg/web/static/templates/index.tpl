<html>
  {{ template "header" }}
  <body>
    <div>My devices</div>
    <table>
     <tr>
       <th>Id</th>
       <th>Name</th>
       <th>Publickey</th>
       <th>Last seen</th>
       <th>Actions</th>
     </tr>

    {{ range $.devices }}
     <tr>
       <td>{{ .id }}</td>
       <td>{{ .name }}</td>
       <td>{{ .pubkey }}</td>
       <td>{{ .lastSeen }}</td>
       <td><a href="/devices/{{ .id }}/download">Download</a></td>
     </tr>
    {{ end }}

   </table>

  </body>
</html>
