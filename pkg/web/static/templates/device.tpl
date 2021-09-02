<html>
  {{ template "header" }}
  <body>
    <br>
    <div class="errors">
    {{ range $.errors }}
    <div>{{ . }}</div>
    {{ end }}
    </div>
    <br>
    <div class="device-tite">New Device:</div>
    <div class="device-form">
      <form action="/new_device" method="POST">
        <label for="name">Name:</label><br>
        <input type="text" id="name" name="name"><br>
        <input type="submit" value="Submit"></input>
      </form>
    </div>
  </body>
</html>
