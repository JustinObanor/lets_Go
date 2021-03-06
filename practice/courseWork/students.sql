PGDMP     2                    w            bookstore_go    11.5    11.5                0    0    ENCODING    ENCODING        SET client_encoding = 'UTF8';
                       false                       0    0 
   STDSTRINGS 
   STDSTRINGS     (   SET standard_conforming_strings = 'on';
                       false                       0    0 
   SEARCHPATH 
   SEARCHPATH     8   SELECT pg_catalog.set_config('search_path', '', false);
                       false                       1262    24700    bookstore_go    DATABASE     �   CREATE DATABASE bookstore_go WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'English_United States.1252' LC_CTYPE = 'English_United States.1252';
    DROP DATABASE bookstore_go;
             postgres    false                        0    0    DATABASE bookstore_go    ACL     1   GRANT ALL ON DATABASE bookstore_go TO justin_go;
                  postgres    false    2847            �            1259    24702    books    TABLE     �   CREATE TABLE public.books (
    studentid character varying NOT NULL,
    firstname character varying NOT NULL,
    lastname character varying NOT NULL,
    classcode numeric(5,2) NOT NULL,
    roomnumber numeric,
    feestobepaid character varying
);
    DROP TABLE public.books;
       public         postgres    false            �            1259    41119    cleaners    TABLE     �   CREATE TABLE public.cleaners (
    workerid integer NOT NULL,
    firstname character varying,
    lastname character varying,
    post character varying
);
    DROP TABLE public.cleaners;
       public         postgres    false            �            1259    41208 	   dormitory    TABLE     �   CREATE TABLE public.dormitory (
    dormitoryid integer,
    studentid character varying,
    workerid integer,
    gardenerid integer,
    securityid integer
);
    DROP TABLE public.dormitory;
       public         postgres    false            �            1259    41125 	   gardeners    TABLE     �   CREATE TABLE public.gardeners (
    gardenerid integer NOT NULL,
    firstname character varying,
    lastname character varying,
    post character varying
);
    DROP TABLE public.gardeners;
       public         postgres    false            �            1259    41131    guards    TABLE     �   CREATE TABLE public.guards (
    securityid integer NOT NULL,
    firstname character varying,
    lastname character varying,
    post character varying,
    timeofwork character varying
);
    DROP TABLE public.guards;
       public         postgres    false                      0    24702    books 
   TABLE DATA               d   COPY public.books (studentid, firstname, lastname, classcode, roomnumber, feestobepaid) FROM stdin;
    public       postgres    false    196   �                 0    41119    cleaners 
   TABLE DATA               G   COPY public.cleaners (workerid, firstname, lastname, post) FROM stdin;
    public       postgres    false    197   �                 0    41208 	   dormitory 
   TABLE DATA               ]   COPY public.dormitory (dormitoryid, studentid, workerid, gardenerid, securityid) FROM stdin;
    public       postgres    false    200   m                 0    41125 	   gardeners 
   TABLE DATA               J   COPY public.gardeners (gardenerid, firstname, lastname, post) FROM stdin;
    public       postgres    false    198   �                 0    41131    guards 
   TABLE DATA               S   COPY public.guards (securityid, firstname, lastname, post, timeofwork) FROM stdin;
    public       postgres    false    199          �
           2606    24709    books books_pkey 
   CONSTRAINT     U   ALTER TABLE ONLY public.books
    ADD CONSTRAINT books_pkey PRIMARY KEY (studentid);
 :   ALTER TABLE ONLY public.books DROP CONSTRAINT books_pkey;
       public         postgres    false    196            �
           2606    41161    cleaners cleaners_pkey 
   CONSTRAINT     Z   ALTER TABLE ONLY public.cleaners
    ADD CONSTRAINT cleaners_pkey PRIMARY KEY (workerid);
 @   ALTER TABLE ONLY public.cleaners DROP CONSTRAINT cleaners_pkey;
       public         postgres    false    197            �
           2606    41179    gardeners gardeners_pkey 
   CONSTRAINT     ^   ALTER TABLE ONLY public.gardeners
    ADD CONSTRAINT gardeners_pkey PRIMARY KEY (gardenerid);
 B   ALTER TABLE ONLY public.gardeners DROP CONSTRAINT gardeners_pkey;
       public         postgres    false    198            �
           2606    41181    guards guards_pkey 
   CONSTRAINT     X   ALTER TABLE ONLY public.guards
    ADD CONSTRAINT guards_pkey PRIMARY KEY (securityid);
 <   ALTER TABLE ONLY public.guards DROP CONSTRAINT guards_pkey;
       public         postgres    false    199            �
           2606    41224 #   dormitory dormitory_gardenerid_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.dormitory
    ADD CONSTRAINT dormitory_gardenerid_fkey FOREIGN KEY (gardenerid) REFERENCES public.gardeners(gardenerid);
 M   ALTER TABLE ONLY public.dormitory DROP CONSTRAINT dormitory_gardenerid_fkey;
       public       postgres    false    200    198    2709            �
           2606    41229 #   dormitory dormitory_securityid_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.dormitory
    ADD CONSTRAINT dormitory_securityid_fkey FOREIGN KEY (securityid) REFERENCES public.guards(securityid);
 M   ALTER TABLE ONLY public.dormitory DROP CONSTRAINT dormitory_securityid_fkey;
       public       postgres    false    199    200    2711            �
           2606    41214 "   dormitory dormitory_studentid_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.dormitory
    ADD CONSTRAINT dormitory_studentid_fkey FOREIGN KEY (studentid) REFERENCES public.books(studentid);
 L   ALTER TABLE ONLY public.dormitory DROP CONSTRAINT dormitory_studentid_fkey;
       public       postgres    false    196    200    2705            �
           2606    41219 !   dormitory dormitory_workerid_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.dormitory
    ADD CONSTRAINT dormitory_workerid_fkey FOREIGN KEY (workerid) REFERENCES public.cleaners(workerid);
 K   ALTER TABLE ONLY public.dormitory DROP CONSTRAINT dormitory_workerid_fkey;
       public       postgres    false    197    200    2707               �   x�M��N�0���)x��?qh�A�R!Q�.\��*N��nEޞMP�^�2�f����'��	u%%�T|�]���p?�p	x;�hi`�
��[�8R�iB3�=%�J�qDEn"����nvhS�ś�%����+۬F�p]���#�,��z�<�s<y|Ĕ��3w���yp��p},_03��n���R�B�v8S���[(��m->+!��VF         �   x�=��
�0���)���z���Ҟ<	^���6�MJ_�h!�=}3�-��XF�C h�CVȐ�g�D�h���h,�&���2�����@OBs�;<����8��5Քx�>�G�V���/t�j*܅���i�a�Z�J�FG�0[4���N.���J�/�IC�            x������ � �         z   x�3���L�HM��t�IL��t+��+QpO,JI��2��J�M-�t�/-*�I�D�5����+.��1JR9�3SRa�&����%P��%�(��@��9C2�s�9��ä�8}RA&塈��qqq ٱ6�         �   x��˻
�@���W���,"	� 66��!��DؿW!�U��Y"'!�PF\��$�H<{RQ�6�;~�Y��ءCc�?�Fn=n+��q�œ=�4�j�\>wnqz�jq����Ȓ�N'��1��	�m����1�6�X�     