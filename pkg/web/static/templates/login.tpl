<html>
  {{ template "header" }}
  <body>
    {{ range $.errors }}
    <div>{{ . }}</div>
    {{ end }}
    <br>
    <div class="login-tite">Login:</div>
    <div class="login-form">
      <form action="/login" method="POST">
        <label for="email">Email:</label><br>
        <input type="text" id="email" name="email"><br>
        <label for="password">Password:</label><br>
        <input type="text" id="password" name="password">
        <input type="submit" value="Submit"></input>
      </form>
    </div>
  </body>
</html>
