# Run as
#
#   python3 setup.py build_ext --build-lib=src
#
# to build Cython modules during development and place .so files into src/pyhexz.
from setuptools import setup, Extension
from Cython.Build import cythonize
import os

setup(
    name="pyhexz",
    ext_modules=cythonize(
        [Extension("pyhexz.hexc", ["src/pyhexz/hexc.py"])]
        + (
            [  # Building the pyhexz.cclib extension assumes that
                # the C++ libraries in ../cpp were already built.
                Extension(
                    "pyhexz.ccapi",
                    ["src/pyhexz/ccapi.pyx"],
                    libraries=["hexz_pyapi"],
                    library_dirs=["../cpp/build"],
                    include_dirs=["../cpp"],
                    language="c++",
                    extra_compile_args=["-std=c++17"],
                ),
            ]
            if os.getenv("HEXZ_SKIP_BUILD_CCAPI") is None
            else []
        ),
        language_level="3",
    ),
)
