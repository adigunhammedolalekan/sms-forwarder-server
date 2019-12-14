# SMS-forwarder Server
An sms forwarder http API server.

* user creates an account
* user registers for sms listener using email(s) of the listener - listener are those who listens/receives the forwarded SMS
* When mobile client receives SMS, sms forwarder forwards it to server and server forwards it to all registered listener

## Use cases
* In a system whereby multiple people needs to have access to a local mobile SMS