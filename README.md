# logbox
A very simple centralized logging system

# User Interface
You can access the user interface at http://localhost:8080

# Forwarding logs 

You can send it logs on a simple TCP connection to localhost:8888

For example, you can use netcat:

`cat logfile.txt | nc -q 0 localhost 8888`
