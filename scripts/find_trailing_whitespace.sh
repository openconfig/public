#!/bin/bash
#
# MIT License
# 
# Copyright (c) 2020 Harutaka Kawamura
# 
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
# 
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
# 
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

set -e

# Scan all tracked YANG files in a single git grep process.
tw_lines=""
if tw_lines=$(git grep -n -I -E '[[:blank:]]+$' -- 'release/models/*.yang'); then
  :
else
  grep_status=$?
  if [ "$grep_status" -ne 1 ]; then
    exit "$grep_status"
  fi
fi

exit_code=0

# If tw_lines is not empty, change the exit code to 1 to fail the CI.
if [ ! -z "$tw_lines" ]; then
  echo -e "\n***** Lines containing trailing whitespace *****\n"
  echo -e "${tw_lines[@]}"
  echo -e "\nFailed. Try runing the following command:\n"
  echo -e "\nfind release/models -type f -name '*\\.yang' -exec sed -i 's/[[:space:]]\\+$//' {} \\+\n"
  exit_code=1
fi

exit $exit_code
