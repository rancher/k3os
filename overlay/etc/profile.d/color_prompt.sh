# Setup a red prompt for root and a green one for users, on non-dumb terminals.

case "$TERM" in
"dumb")
	;;
xterm*|rxvt*|eterm*|screen*)
	NORMAL="\[\e[0m\]"
	RED="\[\e[1;31m\]"
	GREEN="\[\e[1;32m\]"
	;;
esac

if [ "$(id -u)" = 0 ]; then
	PS1="$RED\h [$NORMAL\w$RED]# $NORMAL"
else
	PS1="$GREEN\h [$NORMAL\w$GREEN]\$ $NORMAL"
fi
