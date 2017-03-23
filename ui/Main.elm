module Main exposing (..)

import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Http
import Json.Decode as Decode


main =
    Html.program
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }


type alias Model =
    { party : String
    }


init : ( Model, Cmd Msg )
init =
    ( Model ""
    , getParties
    )


type Msg
    = MorePlease
    | GotParties (Result Http.Error String)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        MorePlease ->
            ( model, getParties )

        GotParties (Ok party) ->
            ( Model party, Cmd.none )

        GotParties (Err _) ->
            ( model, Cmd.none )


view : Model -> Html Msg
view model =
    div []
        [ button [ onClick MorePlease ] [ text "More Please!" ]
        , br [] []
        , li [] [ text model.party ]
        ]


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none


getParties : Cmd Msg
getParties =
    let
        url =
            "http://127.0.0.1:3000/api/parties"
    in
        Http.send GotParties (Http.get url decodePartyData)


decodePartyData : Decode.Decoder String
decodePartyData =
    Decode.index 0 (Decode.at [ "hash" ] Decode.string)
