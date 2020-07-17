urldecode() { : "${*//+/ }"; echo -e "${_//%/\\x}"; }

urldecode $1
