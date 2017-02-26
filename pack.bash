#!/bin/bash
WORKDIR=$(pwd)
package(){
		set -e
		echo "Using ./pkg directory"
        mkdir -p pkg
        if [ -f HASH ]; then
        	echo "Renaming HASH to HASH.old"
        	mv HASH HASH.old
        fi
        cd bin
        echo "Creating HASH file"
        for i in $(ls | grep -v "VERSION"); do sha256sum $i >> $WORKDIR/HASH; done
        cd $WORKDIR
        echo "Packaging all in ./bin"
        for i in $(ls bin|grep -v "VERSION"); do zip pkg/$i.zip bin/$i README.md LICENSE.md HASH; done
		echo "Done."
		echo ""
}

BINFILES=$(ls bin)
if [ -z "$BINFILES" ]; then
echo "Run 'make cross' first!"
exit 1
fi
if [ ! -d bin ]; then
echo "Run 'make cross' first!"
exit 1
fi
package
