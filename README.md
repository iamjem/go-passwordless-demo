# Passwordless Demo
This is the project source code from my [blog post](#TODO) on [pixeldonor.com](http://www.pixeldonor.com).

# Running the Project On Heroku
If you're a heroku fan, simply use the heroku button and you'll be up and running with the demo app in moments.

# Running the Project Local
Should you wish to run the project local, you'll want to make sure you've set all the necessary environment variables.

# Generating Random Keys
There's two required environment variables `AUTH_KEY` and `ENCRYPT_KEY` which should be random strings. You can use [this snippet](http://play.golang.org/p/TKd3pMLx7c) on play.golang.org to generate your own. Remember, the `ENCRYPT_KEY` should be of a fixed length (16, 24, or 32) so you'll need to adjust the string accordingly.
