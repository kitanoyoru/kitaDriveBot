mods=$(go list -f '{{.Dir}}' -m | xargs)

for mod in $mods; do
    (cd "$mod"; go test ./...)
done
