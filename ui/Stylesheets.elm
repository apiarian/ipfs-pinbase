port module Stylesheets exposing (..)

import Css.File exposing (CssFileStructure, CssCompilerProgram)
import PinCss


port files : CssFileStructure -> Cmd msg


fileStructure : CssFileStructure
fileStructure =
    Css.File.toFileStructure
        [ ( "pinbase.css", Css.File.compile [ PinCss.css ] ) ]


main : CssCompilerProgram
main =
    Css.File.compiler files fileStructure
