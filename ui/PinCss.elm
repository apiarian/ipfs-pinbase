module PinCss exposing (..)

import Css exposing (..)
import Css.Elements exposing (body, li, ul)
import Css.Namespace exposing (namespace)


type CssClasses
    = Main
    | PartyList
    | PartyEntity
    | PartyDescription
    | PartyHash
    | PartyPins
    | PinEntity
    | PinEditing
    | PinHeader
    | PinName
    | PinHash
    | PinAliasesBox
    | PinAliasesLabel
    | PinAliasesList
    | PinAlias
    | PinPinBox
    | PinStatusText
    | PinStatusGood
    | PinStatusPending
    | PinStatusError


pinspace : String
pinspace =
    "pinbase"


css : Stylesheet
css =
    (stylesheet << namespace pinspace)
        [ class Main
            [ maxWidth (px 800)
            , margin auto
            ]
        , class PartyDescription
            [ fontSize xLarge
            , fontWeight bold
            ]
        , class PartyHash
            [ hashStyle
            ]
        , class PinEditing
            [ border3 (px 1) solid (rgb 0x00 0x33 0xFF)
            ]
        , class PinName
            [ fontSize large
            , fontWeight bold
            ]
        , class PinHash
            [ hashStyle
            ]
        , class PinAliasesLabel
            [ marginBottom zero
            ]
        , class PinAliasesBox
            [ children
                [ ul
                    [ marginTop zero
                    ]
                ]
            ]
        , class PinPinBox
            [ float left
            , padding (em 0.5)
            ]
        , class PinStatusText
            [ padding (em 0.5)
            ]
        ]


hashColor : Color
hashColor =
    rgb 0x88 0x88 0x88


hashStyle : Mixin
hashStyle =
    mixin
        [ fontSize small
        , color hashColor
        ]
