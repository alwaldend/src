# def al_go_checkers(
#         name,
#         srcs = [],
#         go_src = GO_SRC,
#         run_args_src = RUN_ARGS_SRC):
#     """
#     Generate golang checkers for
#     """
#     kwargs = {"srcs": [run_args_src], "data": srcs + [go_src]}
#     native.sh_binary(
#         name = "{}-stylua-fix".format(name),
#         args = stylua_args,
#         **stylua_kwargs
#     )
#     native.sh_test(
#         name = "{}-stylua-test".format(name),
#         args = [stylua_args[0], "--check"] + stylua_args[1:],
#         size = "small",
#         **stylua_kwargs
#     )
