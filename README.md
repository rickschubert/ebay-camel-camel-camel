Ebay Camel Camel Camel
=====================

A small web scraper which gives you the links to auctions that are ending soon.

# How to use
This script is not intended to be a CLI, meaning you have to change the values yourself. Simply define in the `main` function the product you want to look for, the maximum price you are willing to pay and how much time you want auctions to have maximum left.

Be aware that this script is written against the ebay UK website - to use it for other regions, adjust the base URL in the `getAuctions()` function and the currency symbol (currently `£`).
