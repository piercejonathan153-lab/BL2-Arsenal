#!/bin/sh
# usage: ./publish.sh path/to/bl2_arsenal_vX.Y.Z.sdkmod "note one" "note two" ...
# copies the .sdkmod into releases/, rewrites latest.json and prepends CHANGELOG.md; then commit + push.
set -e
f="$1"; shift
python3 - "$f" "$@" <<'PY'
import json,hashlib,os,re,shutil,sys
f=sys.argv[1]; notes=sys.argv[2:]
v=re.search(r'(\d+\.\d+\.\d+)',os.path.basename(f)).group(1)
dst=os.path.join('releases',os.path.basename(f)); shutil.copy(f,dst)
lj=json.load(open('latest.json'))
lj.update(version=v,file=os.path.basename(f),size=os.path.getsize(dst),sha256=hashlib.sha256(open(dst,'rb').read()).hexdigest().upper(),notes=notes)
json.dump(lj,open('latest.json','w'),indent=1)
c=open('CHANGELOG.md').read().split('\n',2)
open('CHANGELOG.md','w').write(c[0]+'\n\n## v'+v+'\n'+''.join('- '+n+'\n' for n in notes)+'\n'+c[2].lstrip('\n'))
print('published',v)
PY
