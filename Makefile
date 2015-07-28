# 
# go-sonos
# ========
# 
# Copyright (c) 2012, Ian T. Richards <ianr@panix.com>
# All rights reserved.
# 
# Redistribution and use in source and binary forms, with or without
# modification, are permitted provided that the following conditions
# are met:
# 
#   * Redistributions of source code must retain the above copyright notice,
#     this list of conditions and the following disclaimer.
#   * Redistributions in binary form must reproduce the above copyright
#     notice, this list of conditions and the following disclaimer in the
#     documentation and/or other materials provided with the distribution.
# 
# THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
# "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
# LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
# A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
# HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
# SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
# TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
# PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
# LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
# NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
# SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
# 

GO = go
PACKAGE = github.com/ianr0bkny/go-sonos

all :: examples
	$(GO) install -v
	$(GO) install -v $(PACKAGE)/cscl
	$(GO) install -v $(PACKAGE)/csweb

clean ::
	$(GO) clean -i -x
	$(GO) clean -i -x $(PACKAGE)/cscl
	$(GO) clean -i -x $(PACKAGE)/csweb
	$(GO) clean -i -x $(PACKAGE)/examples/devices
	$(GO) clean -i -x $(PACKAGE)/examples/discovery
	$(GO) clean -i -x $(PACKAGE)/examples/googletv
	rm -rf $(GOPATH)/pkg/*/$(PACKAGE)

wc ::
	wc -l *.go */*.go

longlines ::
	egrep '.{120,}' *.go */*.go

coverage ::
	$(GO) test -test.run Coverage

discovery ::
	$(GO) test -test.run Discovery

examples ::
	$(GO) install -v $(PACKAGE)/examples/browse
	$(GO) install -v $(PACKAGE)/examples/composers
	$(GO) install -v $(PACKAGE)/examples/devices
	$(GO) install -v $(PACKAGE)/examples/discovery
	$(GO) install -v $(PACKAGE)/examples/googletv
