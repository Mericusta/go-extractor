# go-extractor

## Features

- MakeUnitTest
  - generate unit test func
  - if func is a generic func, then the number of 'typeArgs' must greater than the number of type params

## Todo List

- GoProjectMeta
  - struct
    - GoProjectMeta
  - method
    - SearchPackageMeta
  - func
    - ExtractGoProjectMeta
- GoPackageMeta
  - struct
    - GoPackageMeta
  - method
    - extractVar
    - extractFunc
    - extractStruct
    - extractMethod
    - extractInterface
    - SearchFileMeta
    - SearchVarMeta
    - SearchFuncMeta
    - SearchStructMeta
    - SearchInterfaceMeta
    - StructNames
    - InterfaceNames
    - FunctionNames
  - func
    - newGoPackageMeta
    - ExtractGoPackageMeta
    - ExtractGoPackageMetaWithSpecPaths
    - extractGoPackageMeta
    - ExtractAll
- GoFileMeta
- GoVarMeta
- GoFuncMeta
  - struct
    - GoFuncMeta
  - func
    - Extract
      - from package
      - from file
    - newGoFuncMeta
