Description: Generating two step code and copying it to clipboard for convenient

Platform: Linux (I only test on ubuntu)

Prerequirement: sudo apt-get install xsel

Usage:

You could download the executable file totp and add it to your system PATH.

There is three paramters:

--key,  set the key, if not, it will use the default one place in ~/.totp/key

--len,  set the length of the output code, default is 6

--step, set the step, default is 30
