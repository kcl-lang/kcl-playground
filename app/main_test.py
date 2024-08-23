import requests
from main import website, router


def test_website():
    assert "8000" in website.url()


def test_router():
    assert "8001" in router.url()
