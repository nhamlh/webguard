<html>
  {{ template "header" }}
  <body>
    <div>My devices</div>
    <table>
     <tr>
       <th>Name</th>
       <th>Publickey</th>
       <th>Last seen</th>
     </tr>

    {{ range $.devices }}
     <tr>
       <td>{{ .name }}</td>
       <td>{{ .pubkey }}</td>
       <td>{{ .lastSeen }}</td>
     </tr>
    {{ end }}

   </table>

  </body>
</html>
