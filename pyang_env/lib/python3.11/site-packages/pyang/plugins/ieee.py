"""IEEE usage guidelines plugin
See http://standards.ieee.org/develop/regauth/tut/ieeeurn.pdf
"""

import optparse
import re

from pyang import plugin
from pyang import statements
from pyang import error
from pyang.error import err_add
from pyang.plugins import lint

def pyang_plugin_init():
    plugin.register_plugin(IEEEPlugin())

class IEEEPlugin(lint.LintPlugin):
    def __init__(self):
        lint.LintPlugin.__init__(self)
        self.modulename_prefixes = ['ieee']

    def add_opts(self, optparser):
        optlist = [
            optparse.make_option("--ieee",
                                 dest="ieee",
                                 action="store_true",
                                 help="Validate the module(s) according to " \
                                 "IEEE rules."),
            ]
        optparser.add_options(optlist)

    def setup_ctx(self, ctx):
        if not ctx.opts.ieee:
            return
        self._setup_ctx(ctx)

        error.add_error_code(
           'IEEE_BAD_NAMESPACE_VALUE', 4,
           'the namespace should be on the form '
           'urn:ieee:std:{IEEE standard designation}:yang:%s')

        statements.add_validation_fun(
            'grammar', ['namespace'],
            lambda ctx, s: self.v_chk_namespace(ctx, s))

    def v_chk_namespace(self, ctx, stmt):
        r = 'urn:ieee:std:.*:yang:' + stmt.i_module.arg
        if re.match(r, stmt.arg) is None:
            err_add(ctx.errors, stmt.pos, 'IEEE_BAD_NAMESPACE_VALUE',
                    stmt.i_module.arg)
