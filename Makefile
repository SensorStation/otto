PLUGINS := $(wildcard plugins/*)
SUBDIRS = $(PLUGINS) mock otto

all: $(SUBDIRS)

$(SUBDIRS):
	$(MAKE) -C $@

.PHONY: all $(SUBDIRS)
