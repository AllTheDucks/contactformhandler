## Contact Form Handler
This is a basic service which listens for requests over http, then emails the
posted data to a list of predefined email addresses.

An example of how to run the compiled command.

```sh
contactformhandler -from="Contact Form <contactform@mycompany.com>" -to="Me <my.email@address.com>" -smtpuser=mysmtpusername -smtppass="mysecretsmtppass" -smtpserver="my.smtp.server.com" -corsurl="http://www.mycompany.com"
```