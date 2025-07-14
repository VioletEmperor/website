#kill "$(ps -e | grep docker-compose | awk '{print $1}')"
kill $(pgrep docker-compose)
